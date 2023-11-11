package api

import (
	"bar/autogen"
	"bar/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) createCashMovement(c mongo.SessionContext, createdBy uuid.UUID, createdByName string, toAdd int64, reason string) error {
	latest, err := s.DBackend.GetLatestCashMovement(c)
	if err != nil {
		return err
	}

	newCashMovement := &models.CashMovement{
		CashMovement: autogen.CashMovement{
			Amount:        latest.Amount + toAdd,
			OldAmount:     latest.Amount,
			CreatedAt:     uint64(time.Now().Unix()),
			CreatedBy:     createdBy,
			CreatedByName: createdByName,
			Reason:        reason,
			Id:            uuid.New(),
		},
	}

	if err := s.DBackend.CreateCashMovement(c, newCashMovement); err != nil {
		return err
	}

	return nil
}

// (GET /cash_movements)
func (s *Server) GetCashMovements(c echo.Context, params autogen.GetCashMovementsParams) error {
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	search := ""
	if params.Search != nil {
		search = *params.Search
	}

	count, err := s.DBackend.CountAllCashMovements(c.Request().Context(), search)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetAllCashMovements(c.Request().Context(), dbpage, limit, search)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	cashMovements := make([]autogen.CashMovement, len(data))
	for i, r := range data {
		cashMovements[i] = r.CashMovement
	}

	autogen.GetCashMovements200JSONResponse{
		CashMovements: cashMovements,
		Limit:         limit,
		Page:          page,
		MaxPage:       maxPage,
	}.VisitGetCashMovementsResponse(c.Response())
	return nil
}

// (POST /cash_movements)
func (s *Server) CreateCashMovement(ctx echo.Context) error {
	return nil
}

// (DELETE /cash_movements/{cash_movement_id})
func (s *Server) DeleteCashMovement(ctx echo.Context, cashMovementId autogen.UUID) error {
	return nil
}

package api

import (
	"bar/autogen"
	"bar/internal/models"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (GET /refills)
func (s *Server) GetRefills(c echo.Context, params autogen.GetRefillsParams) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

	var startsAt uint64 = 0
	if params.StartDate != nil {
		startsAt = uint64(params.StartDate.Unix())
	}
	var endsAt uint64 = math.MaxInt64
	if params.EndDate != nil {
		endsAt = uint64(params.EndDate.Unix())
	}

	count, err := s.DBackend.CountAllRefills(c.Request().Context(), startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var page uint64 = 0
	if params.Page != nil {
		page = *params.Page
	}

	var size uint64 = 10
	if params.Limit != nil {
		size = *params.Limit
	}

	if size > 100 {
		size = 100
	}

	maxPage := uint64(count) / size
	if page > maxPage {
		page = maxPage
	}

	data, err := s.DBackend.GetAllRefills(c.Request().Context(), page, size, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var refills []autogen.Refill

	for _, refill := range data {
		refills = append(refills, refill.Refill)
	}

	logrus.Infof("Refills have been retrieved by %s", adminID)
	autogen.GetRefills200JSONResponse(refills).VisitGetRefillsResponse(c.Response())
	return nil
}

// (GET /account/refills)
func (s *Server) GetSelfRefills(c echo.Context, params autogen.GetSelfRefillsParams) error {
	// Get account from cookie
	logged := c.Get("userLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	accountID := c.Get("userAccountID").(string)

	var startsAt uint64 = 0
	if params.StartDate != nil {
		startsAt = uint64(params.StartDate.Unix())
	}

	var endsAt uint64 = math.MaxInt64
	if params.EndDate != nil {
		endsAt = uint64(params.EndDate.Unix())
	}

	count, err := s.DBackend.CountRefills(c.Request().Context(), accountID, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var page uint64 = 0
	if params.Page != nil {
		page = *params.Page
	}

	var size uint64 = 10
	if params.Limit != nil {
		size = *params.Limit
	}

	if size > 100 {
		size = 100
	}

	maxPage := uint64(count) / size
	if page > maxPage {
		page = maxPage
	}

	data, err := s.DBackend.GetRefills(c.Request().Context(), accountID, page, size, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var refills []autogen.Refill

	for _, refill := range data {
		refills = append(refills, refill.Refill)
	}

	autogen.GetSelfRefills200JSONResponse{
		Refills: &refills,
		Limit:   size,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetSelfRefillsResponse(c.Response())
	return nil
}

// (GET /accounts/{account_id}/refills)
func (s *Server) GetAccountRefills(c echo.Context, accountId autogen.UUID, params autogen.GetAccountRefillsParams) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

	var startsAt uint64 = 0
	if params.StartDate != nil {
		startsAt = uint64(params.StartDate.Unix())
	}

	var endsAt uint64 = math.MaxInt64
	if params.EndDate != nil {
		endsAt = uint64(params.EndDate.Unix())
	}

	count, err := s.DBackend.CountRefills(c.Request().Context(), accountId.String(), startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var page uint64 = 0
	if params.Page != nil {
		page = *params.Page
	}

	var size uint64 = 10
	if params.Limit != nil {
		size = *params.Limit
	}

	if size > 100 {
		size = 100
	}

	maxPage := uint64(count) / size
	if page > maxPage {
		page = maxPage
	}

	data, err := s.DBackend.GetRefills(c.Request().Context(), accountId.String(), page, size, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var refills []autogen.Refill

	for _, refill := range data {
		refills = append(refills, refill.Refill)
	}

	logrus.Infof("Refills have been retrieved by %s", adminID)
	autogen.GetAccountRefills200JSONResponse{
		Refills: &refills,
		Limit:   size,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetAccountRefillsResponse(c.Response())
	return nil
}

// (POST /accounts/{account_id}/refills)
func (s *Server) PostRefill(c echo.Context, accountId autogen.UUID, params autogen.PostRefillParams) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

	account, err := s.DBackend.GetAccount(c.Request().Context(), accountId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	refill := &models.Refill{
		Refill: autogen.Refill{
			AccountId: accountId,
			Amount:    params.Amount,
			Id:        uuid.New(),
			IssuedAt:  uint64(time.Now().Unix()),
			IssuedBy:  uuid.MustParse(adminID),
			State:     autogen.Valid,
		},
	}

	account.Balance += int64(params.Amount)

	_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		err := s.DBackend.CreateRefill(c.Request().Context(), refill)
		if err != nil {
			return nil, errors.New("failed to create refill")
		}

		err = s.DBackend.UpdateAccount(c.Request().Context(), account)
		if err != nil {
			return nil, errors.New("failed to update account")
		}
		return nil, nil
	})
	if err != nil {
		return Error500(c)
	}

	logrus.Infof("Refill %s has been created by %s", refill.Id, adminID)
	autogen.PostRefill201JSONResponse(refill.Refill).VisitPostRefillResponse(c.Response())
	return nil
}

// (PATCH /accounts/{account_id}/refills/{refill_id})
func (s *Server) PatchRefillId(c echo.Context, accountId autogen.UUID, refillId autogen.UUID, params autogen.PatchRefillIdParams) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

	account, err := s.DBackend.GetAccount(c.Request().Context(), accountId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	refill, err := s.DBackend.GetRefill(c.Request().Context(), refillId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorRefillNotFound(c)
		}
		return Error500(c)
	}

	oldState := refill.State
	refill.State = params.State

	if oldState == autogen.Valid && params.State == autogen.Canceled {
		account.Balance -= int64(refill.Amount)
	} else if oldState == autogen.Canceled && params.State == autogen.Valid {
		account.Balance += int64(refill.Amount)
	}

	_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		err := s.DBackend.UpdateRefill(ctx, refill)
		if err != nil {
			return nil, errors.New("failed to update refill")
		}

		err = s.DBackend.UpdateAccount(ctx, account)
		if err != nil {
			return nil, errors.New("failed to update account")
		}

		return nil, nil
	})
	if err != nil {
		return Error500(c)
	}

	logrus.Infof("Refill %s has been updated by %s", refill.Id, adminID)
	autogen.PatchRefillId200JSONResponse(refill.Refill).VisitPatchRefillIdResponse(c.Response())
	return nil
}

// (DELETE /accounts/{account_id}/refills/{refill_id})
func (s *Server) MarkDeleteRefill(c echo.Context, accountId autogen.UUID, refillId autogen.UUID) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

	_, err := s.DBackend.GetRefill(c.Request().Context(), refillId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorRefillNotFound(c)
		}
		return Error500(c)
	}

	err = s.DBackend.MarkDeleteRefill(c.Request().Context(), refillId.String(), adminID)
	if err != nil {
		return Error500(c)
	}

	logrus.Infof("Refill %s has been marked deleted by %s", refillId, adminID)
	autogen.MarkDeleteAccountId204Response{}.VisitMarkDeleteAccountIdResponse(c.Response())
	return nil
}

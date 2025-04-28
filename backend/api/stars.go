package api

import (
	"bar/autogen"
	"bar/internal/models"
	"errors"
	"math"
	"time"

	//	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (POST /accounts/{account_id}/stars)
func (s *Server) PostStarring(c echo.Context, accountId string, params autogen.PostStarringParams) error {
	admin, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	account, err := s.DBackend.GetAccountByCard(c.Request().Context(), accountId)
	if account == nil {
		account, err = s.DBackend.GetAccount(c.Request().Context(), accountId)
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	starring := &models.Starring{
		Starring: autogen.Starring{
			AccountId:    account.Id,
			AccountName:  account.Name(),
			Amount:       params.Amount,
			Type:         params.Type,
			Id:           uuid.New(),
			IssuedAt:     uint64(time.Now().Unix()),
			IssuedBy:     admin.Id,
			IssuedByName: admin.Name(),
			State:        autogen.StarringStateValid,
		},
	}

	account.Points += int64(params.Amount)

	_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		err := s.DBackend.CreateStarring(ctx, starring)
		if err != nil {
			return nil, errors.New("failed to create starring")
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

	logrus.WithField("starring", starring.Id.String()).WithField("account", account.Name()).Info("Starring created")
	autogen.PostStarring201JSONResponse(starring.Starring).VisitPostStarringResponse(c.Response())
	return nil
}

// (GET /stars)
func (s *Server) GetStarrings(c echo.Context, params autogen.GetStarringsParams) error {
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	logrus.Debug("GetStarrings, called with startsAt: ", *params.StartDate, " endsAt: ", *params.EndDate, " name: ", *params.Name, " page: ", *params.Page, " size: ", *params.Limit)

	var name string
	if params.Name != nil {
		name = string(*params.Name)
	}
	var startsAt uint64 = 0
	if params.StartDate != nil {
		t, err := time.Parse("2006-01-02", *params.StartDate) // 2006-01-02 is the reference time in Go
		if err == nil {
			startsAt = uint64(t.Unix())
		}
	}
	var endsAt uint64 = math.MaxInt64
	if params.EndDate != nil {
		t, err := time.Parse("2006-01-02", *params.EndDate) // Putting the same date doesn't activate the date filter
		if err == nil {
			endsAt = uint64(t.Unix())
		}
	}

	logrus.Debug("GetStarrings, query with startsAt: ", startsAt, " endsAt: ", endsAt, " name: ", name, "")

	count, err := s.DBackend.CountAllStarrings(c.Request().Context(), name, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetAllStarrings(c.Request().Context(), dbpage, limit, name, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var starrings []autogen.Starring

	for _, starring := range data {
		starrings = append(starrings, starring.Starring)
	}

	autogen.GetStarrings200JSONResponse{
		Stars:   starrings,
		Limit:   limit,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetStarringsResponse(c.Response())
	return nil
}

// (GET /account/stars)
func (s *Server) GetSelfStarring(c echo.Context, params autogen.GetSelfStarringParams) error {
	// Get account from cookie
	_, err := MustGetUser(c)
	if err != nil {
		return nil
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

	count, err := s.DBackend.CountStarrings(c.Request().Context(), accountID, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetStarrings(c.Request().Context(), accountID, dbpage, limit, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var starrings []autogen.Starring

	for _, starring := range data {
		starrings = append(starrings, starring.Starring)
	}

	autogen.GetSelfStarring200JSONResponse{
		Stars:   starrings,
		Limit:   limit,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetSelfStarringResponse(c.Response())
	return nil
}

// (GET /accounts/{account_id}/stars)
func (s *Server) GetAccountStarring(c echo.Context, accountId string, params autogen.GetAccountStarringParams) error {
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	account, err := s.DBackend.GetAccountByCard(c.Request().Context(), accountId)
	if account == nil {
		account, err = s.DBackend.GetAccount(c.Request().Context(), accountId)
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	var startsAt uint64 = 0
	if params.StartDate != nil {
		startsAt = uint64(params.StartDate.Unix())
	}

	var endsAt uint64 = math.MaxInt64
	if params.EndDate != nil {
		endsAt = uint64(params.EndDate.Unix())
	}

	count, err := s.DBackend.CountStarrings(c.Request().Context(), account.Id.String(), startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetStarrings(c.Request().Context(), account.Id.String(), dbpage, limit, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var starrings []autogen.Starring

	for _, starring := range data {
		starrings = append(starrings, starring.Starring)
	}

	autogen.GetAccountStarring200JSONResponse{
		Stars:   starrings,
		Limit:   limit,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetAccountStarringResponse(c.Response())
	return nil
}

// (PATCH /accounts/{account_id}/stars/{starring_id})
func (s *Server) PatchStarringId(c echo.Context, accountId autogen.UUID, starringId autogen.UUID, params autogen.PatchStarringIdParams) error {
	admin, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	account, err := s.DBackend.GetAccount(c.Request().Context(), accountId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	starring, err := s.DBackend.GetStarring(c.Request().Context(), starringId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorStarringNotFound(c)
		}
		return Error500(c)
	}

	_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		if params.State != nil {
			oldState := starring.State
			starring.State = *params.State

			if oldState == autogen.StarringStateValid && *params.State == autogen.StarringStateCanceled {
				account.Points -= int64(starring.Amount)

				name := admin.Name()

				starring.CanceledBy = &admin.Id
				starring.CanceledByName = &name

			} else if oldState == autogen.StarringStateCanceled && *params.State == autogen.StarringStateValid {
				account.Points += int64(starring.Amount)
				starring.CanceledBy = nil
				starring.CanceledByName = nil

			}
		}

		if params.Type != nil {
			starring.Type = *params.Type
		}

		err := s.DBackend.UpdateStarring(ctx, starring)
		if err != nil {
			return nil, errors.New("failed to update starring")
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

	logrus.WithField("starring", starring.Id.String()).WithField("account", account.Name()).Info("Starring updated")
	autogen.PatchStarringId200JSONResponse(starring.Starring).VisitPatchStarringIdResponse(c.Response())
	return nil
}

// (DELETE /accounts/{account_id}/stars/{starring_id})
func (s *Server) MarkDeleteStarring(c echo.Context, accountId autogen.UUID, starringId autogen.UUID) error {
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	_, err = s.DBackend.GetStarring(c.Request().Context(), starringId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorStarringNotFound(c)
		}
		return Error500(c)
	}

	err = s.DBackend.MarkDeleteStarring(c.Request().Context(), starringId.String(), account.Id.String())
	if err != nil {
		return Error500(c)
	}

	logrus.WithField("starring", starringId.String()).WithField("account", account.Name()).Info("Starring marked for deletion")
	autogen.MarkDeleteAccountId204Response{}.VisitMarkDeleteAccountIdResponse(c.Response())
	return nil
}

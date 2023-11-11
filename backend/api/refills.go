package api

import (
	"bar/autogen"
	"bar/internal/models"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (GET /refills)
func (s *Server) GetRefills(c echo.Context, params autogen.GetRefillsParams) error {
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	var startsAt uint64 = 0
	if params.StartDate != nil {
		startsAt, _ = strconv.ParseUint(*params.StartDate, 10, 64)
	}
	var endsAt uint64 = math.MaxInt64
	if params.EndDate != nil {
		endsAt, _ = strconv.ParseUint(*params.EndDate, 10, 64)
	}

	count, err := s.DBackend.CountAllRefills(c.Request().Context(), startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetAllRefills(c.Request().Context(), dbpage, limit, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var refills []autogen.Refill

	for _, refill := range data {
		refills = append(refills, refill.Refill)
	}

	logrus.Infof("Refills have been retrieved by %s", account.Id)

	autogen.GetRefills200JSONResponse{
		Refills: refills,
		Limit:   limit,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetRefillsResponse(c.Response())
	return nil
}

// (GET /account/refills)
func (s *Server) GetSelfRefills(c echo.Context, params autogen.GetSelfRefillsParams) error {
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

	count, err := s.DBackend.CountRefills(c.Request().Context(), accountID, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetRefills(c.Request().Context(), accountID, dbpage, limit, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var refills []autogen.Refill

	for _, refill := range data {
		refills = append(refills, refill.Refill)
	}

	autogen.GetSelfRefills200JSONResponse{
		Refills: refills,
		Limit:   limit,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetSelfRefillsResponse(c.Response())
	return nil
}

// (GET /accounts/{account_id}/refills)
func (s *Server) GetAccountRefills(c echo.Context, accountId string, params autogen.GetAccountRefillsParams) error {
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

	count, err := s.DBackend.CountRefills(c.Request().Context(), account.Id.String(), startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetRefills(c.Request().Context(), account.Id.String(), dbpage, limit, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var refills []autogen.Refill

	for _, refill := range data {
		refills = append(refills, refill.Refill)
	}

	logrus.Infof("Refills have been retrieved by %s", account.Id)
	autogen.GetAccountRefills200JSONResponse{
		Refills: refills,
		Limit:   limit,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetAccountRefillsResponse(c.Response())
	return nil
}

// (POST /accounts/{account_id}/refills)
func (s *Server) PostRefill(c echo.Context, accountId string, params autogen.PostRefillParams) error {
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

	refill := &models.Refill{
		Refill: autogen.Refill{
			AccountId:    account.Id,
			AccountName:  account.Name(),
			Amount:       params.Amount,
			Type:         params.Type,
			Id:           uuid.New(),
			IssuedAt:     uint64(time.Now().Unix()),
			IssuedBy:     admin.Id,
			IssuedByName: admin.Name(),
			State:        autogen.Valid,
		},
	}

	account.Balance += int64(params.Amount)

	_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		err := s.DBackend.CreateRefill(ctx, refill)
		if err != nil {
			return nil, errors.New("failed to create refill")
		}

		err = s.DBackend.UpdateAccount(ctx, account)
		if err != nil {
			return nil, errors.New("failed to update account")
		}

		if params.Type == autogen.RefillCash {
			err = s.createCashMovement(ctx, admin.Id, admin.Name(), int64(params.Amount), fmt.Sprintf("Recharge %s sur le compte de %s par %s", refill.Id, account.Name(), admin.Name()))
			if err != nil {
				return nil, errors.New("failed to create cash movement")
			}
		}

		return nil, nil
	})
	if err != nil {
		return Error500(c)
	}

	logrus.Infof("Refill %s has been created by %s", refill.Id, account.Id)
	autogen.PostRefill201JSONResponse(refill.Refill).VisitPostRefillResponse(c.Response())
	return nil
}

// (PATCH /accounts/{account_id}/refills/{refill_id})
func (s *Server) PatchRefillId(c echo.Context, accountId autogen.UUID, refillId autogen.UUID, params autogen.PatchRefillIdParams) error {
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

	refill, err := s.DBackend.GetRefill(c.Request().Context(), refillId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorRefillNotFound(c)
		}
		return Error500(c)
	}

	_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		if params.State != nil {
			oldState := refill.State
			refill.State = *params.State

			if oldState == autogen.Valid && *params.State == autogen.Canceled {
				account.Balance -= int64(refill.Amount)

				name := admin.Name()

				refill.CanceledBy = &admin.Id
				refill.CanceledByName = &name

				err = s.createCashMovement(ctx, admin.Id, admin.Name(), -int64(refill.Amount), fmt.Sprintf("Annulation de recharge %s sur le compte de %s par %s", refill.Id, account.Name(), admin.Name()))
				if err != nil {
					return nil, errors.New("failed to create cash movement")
				}
			} else if oldState == autogen.Canceled && *params.State == autogen.Valid {
				account.Balance += int64(refill.Amount)
				refill.CanceledBy = nil
				refill.CanceledByName = nil

				err = s.createCashMovement(ctx, admin.Id, admin.Name(), int64(refill.Amount), fmt.Sprintf("Revalidation de recharge %s sur le compte de %s par %s", refill.Id, account.Name(), admin.Name()))
				if err != nil {
					return nil, errors.New("failed to create cash movement")
				}
			}
		}

		if params.Type != nil {
			refill.Type = *params.Type
		}

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

	logrus.Infof("Refill %s has been updated by %s", refill.Id, account.Id)
	autogen.PatchRefillId200JSONResponse(refill.Refill).VisitPatchRefillIdResponse(c.Response())
	return nil
}

// (DELETE /accounts/{account_id}/refills/{refill_id})
func (s *Server) MarkDeleteRefill(c echo.Context, accountId autogen.UUID, refillId autogen.UUID) error {
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	_, err = s.DBackend.GetRefill(c.Request().Context(), refillId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorRefillNotFound(c)
		}
		return Error500(c)
	}

	err = s.DBackend.MarkDeleteRefill(c.Request().Context(), refillId.String(), account.Id.String())
	if err != nil {
		return Error500(c)
	}

	logrus.Infof("Refill %s has been marked deleted by %s", refillId, account.Id)
	autogen.MarkDeleteAccountId204Response{}.VisitMarkDeleteAccountIdResponse(c.Response())
	return nil
}

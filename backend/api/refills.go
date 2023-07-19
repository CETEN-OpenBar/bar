package api

import (
	"bar/autogen"
	"bar/internal/models"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (GET /refills)
func (s *Server) GetRefills(c echo.Context, params autogen.GetRefillsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	accountID, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	var startsAt int64 = 0
	if params.StartDate != nil {
		startsAt = params.StartDate.Unix()
	}
	var endsAt int64 = math.MaxInt64
	if params.EndDate != nil {
		endsAt = params.EndDate.Unix()
	}

	count, err := s.DBackend.CountAllRefills(startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	page := 0
	if params.Page != nil {
		page = *params.Page
	}

	size := 10
	if params.Limit != nil {
		size = *params.Limit
	}

	if size < 0 {
		size = 10
	}

	if size > 100 {
		size = 100
	}

	maxPage := int(count) / size
	if page > maxPage {
		page = maxPage
	}

	if page < 0 {
		page = 0
	}

	data, err := s.DBackend.GetAllRefills(page, size, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var refills []autogen.Refill

	for _, refill := range data {
		refills = append(refills, refill.Refill)
	}

	logrus.Infof("Refills have been retrieved by %s", accountID)
	autogen.GetRefills200JSONResponse(refills).VisitGetRefillsResponse(c.Response())
	return nil
}

// (GET /account/refills)
func (s *Server) GetSelfRefills(c echo.Context, params autogen.GetSelfRefillsParams) error {
	// Get account from cookie
	sess := s.getUserSess(c)
	accountID, ok := sess.Values["account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(accountID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	var startsAt int64 = 0
	if params.StartDate != nil {
		startsAt = params.StartDate.Unix()
	}

	var endsAt int64 = math.MaxInt64
	if params.EndDate != nil {
		endsAt = params.EndDate.Unix()
	}

	count, err := s.DBackend.CountRefills(accountID, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	page := 0
	if params.Page != nil {
		page = *params.Page
	}

	size := 10
	if params.Limit != nil {
		size = *params.Limit
	}

	if size < 0 {
		size = 10
	}

	if size > 100 {
		size = 100
	}

	maxPage := int(count) / size
	if page > maxPage {
		page = maxPage
	}

	if page < 0 {
		page = 0
	}

	data, err := s.DBackend.GetRefills(accountID, page, size, startsAt, endsAt)
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
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	var startsAt int64 = 0
	if params.StartDate != nil {
		startsAt = params.StartDate.Unix()
	}

	var endsAt int64 = math.MaxInt64
	if params.EndDate != nil {
		endsAt = params.EndDate.Unix()
	}

	count, err := s.DBackend.CountRefills(accountId.String(), startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	page := 0
	if params.Page != nil {
		page = *params.Page
	}

	size := 10
	if params.Limit != nil {
		size = *params.Limit
	}

	if size < 0 {
		size = 10
	}

	if size > 100 {
		size = 100
	}

	maxPage := int(count) / size
	if page > maxPage {
		page = maxPage
	}

	if page < 0 {
		page = 0
	}

	data, err := s.DBackend.GetRefills(accountId.String(), page, size, startsAt, endsAt)
	if err != nil {
		return Error500(c)
	}

	var refills []autogen.Refill

	for _, refill := range data {
		refills = append(refills, refill.Refill)
	}

	logrus.Infof("Refills have been retrieved by %s", adminId)
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
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	refill := &models.Refill{
		Refill: autogen.Refill{
			AccountId: accountId,
			Amount:    params.Amount,
			Id:        uuid.New(),
			IssuedAt:  time.Now().Unix(),
			IssuedBy:  uuid.MustParse(adminId),
			State:     autogen.Valid,
		},
	}

	err := s.DBackend.CreateRefill(refill)
	if err != nil {
		return Error500(c)
	}

	logrus.Infof("Refill %s has been created by %s", refill.Id, adminId)
	autogen.PostRefill201JSONResponse(refill.Refill).VisitPostRefillResponse(c.Response())
	return nil
}

// (DELETE /accounts/{account_id}/refills/{refill_id})
func (s *Server) MarkDeleteRefill(c echo.Context, accountId autogen.UUID, refillId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	_, err := s.DBackend.GetRefill(refillId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorRefillNotFound(c)
		}
		return Error500(c)
	}

	err = s.DBackend.MarkDeleteRefill(refillId.String(), adminId)
	if err != nil {
		return Error500(c)
	}

	logrus.Infof("Refill %s has been marked deleted by %s", refillId, adminId)
	autogen.MarkDeleteAccountId204Response{}.VisitMarkDeleteAccountIdResponse(c.Response())
	return nil
}

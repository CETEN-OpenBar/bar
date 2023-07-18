package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// (GET /refills)
func (s *Server) GetRefills(c echo.Context, params autogen.GetRefillsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (GET /account/refills)
func (s *Server) GetSelfRefills(c echo.Context, params autogen.GetSelfRefillsParams) error {
	// Get account from cookie
	sess := s.getUserSess(c)
	accountID, ok := sess.Values["account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// Get account from database
	account, err := s.DBackend.GetAccount(accountID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	// TODO: implement
	return nil
}

// (GET /accounts/{account_id}/refills)
func (s *Server) GetAccountRefills(c echo.Context, accountId autogen.UUID, params autogen.GetAccountRefillsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (POST /accounts/{account_id}/refills)
func (s *Server) PostRefill(c echo.Context, accountId autogen.UUID, params autogen.PostRefillParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (DELETE /accounts/{account_id}/refills/{refill_id})
func (s *Server) DeleteRefill(c echo.Context, accountId autogen.UUID, refillId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

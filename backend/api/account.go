package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// (GET /account)
func (s *Server) GetAccount(c echo.Context) error {
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

	// Return account
	resp := autogen.GetAccount200JSONResponse{
		Account: &account.Account,
	}
	resp.VisitGetAccountResponse(c.Response())
	return nil
}

// (GET /accounts)
func (s *Server) GetAccounts(c echo.Context, params autogen.GetAccountsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// Set up parameters
	var page int = 1
	var limit int = 10
	if params.Page != nil {
		page = *params.Page
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	page = page - 1
	if page < 0 {
		page = 0
	}
	if limit < 0 {
		limit = 0
	}
	if limit > 100 {
		limit = 100
	}

	// Calculate max page
	count, err := s.DBackend.CountAccounts()
	if err != nil {
		return Error500(c)
	}

	maxPage := int(count) / limit

	// Get accounts from database
	accounts, err := s.DBackend.GetAccounts(page, limit)
	if err != nil {
		return Error500(c)
	}

	var ac []autogen.Account
	for _, account := range accounts {
		ac = append(ac, account.Account)
	}

	autogen.GetAccounts200JSONResponse{
		Accounts: ac,
		Limit:    limit,
		Page:     page,
		MaxPage:  maxPage,
	}.VisitGetAccountsResponse(c.Response())
	return nil
}

// (POST /accounts)
func (s *Server) PostAccounts(c echo.Context) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (DELETE /accounts/{account_id})
func (s *Server) DeleteAccountId(c echo.Context, accountId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (GET /accounts/{account_id})
func (s *Server) GetAccountId(c echo.Context, accountId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /accounts/{account_id})
func (s *Server) PatchAccountId(c echo.Context, accountId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

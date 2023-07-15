package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
)

// (GET /account)
func (s *Server) GetAccount(c echo.Context) error {
	// TODO: implement
	return nil
}

// (POST /account/transactions)
func (s *Server) PostTransactions(c echo.Context) error {
	// TODO: implement
	return nil
}

// (GET /accounts)
func (s *Server) GetAccounts(c echo.Context, params autogen.GetAccountsParams) error {
	// TODO: implement
	return nil
}

// (POST /accounts)
func (s *Server) PostAccounts(c echo.Context) error {
	// TODO: implement
	return nil
}

// (DELETE /accounts/{account_id})
func (s *Server) DeleteAccountId(c echo.Context, accountId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (GET /accounts/{account_id})
func (s *Server) GetAccountId(c echo.Context, accountId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (PATCH /accounts/{account_id})
func (s *Server) PatchAccountId(c echo.Context, accountId autogen.UUID) error {
	// TODO: implement
	return nil
}

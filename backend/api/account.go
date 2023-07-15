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
		resp := autogen.GetAccount401JSONResponse{
			Message:   autogen.MsgAccountNotFound,
			ErrorCode: autogen.ErrAccountNotFound,
		}
		resp.VisitGetAccountResponse(c.Response())
		return nil
	}

	// Get account from database
	account, err := s.DBackend.GetAccount(accountID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())

			resp := autogen.GetAccount401JSONResponse{
				Message:   autogen.MsgAccountNotFound,
				ErrorCode: autogen.ErrAccountNotFound,
			}
			resp.VisitGetAccountResponse(c.Response())
			return nil
		}
		resp := autogen.GetAccount500JSONResponse{
			Message:   autogen.MsgInternalServerError,
			ErrorCode: autogen.ErrInternalServerError,
		}
		resp.VisitGetAccountResponse(c.Response())
		return nil
	}

	// Return account
	resp := autogen.GetAccount200JSONResponse{
		Account: &account.Account,
	}
	resp.VisitGetAccountResponse(c.Response())
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

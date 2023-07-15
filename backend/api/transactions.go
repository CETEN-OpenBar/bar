package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
)

// (GET /accounts/{account_id}/transactions)
func (s *Server) GetAccountTransactions(c echo.Context, accountId autogen.UUID, params autogen.GetAccountTransactionsParams) error {
	// TODO: implement
	return nil
}

// (GET /account/transactions)
func (s *Server) GetCurrentAccountTransactions(c echo.Context, params autogen.GetCurrentAccountTransactionsParams) error {
	// TODO: implement
	return nil
}

// (DELETE /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) DeleteTransactionId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (GET /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) GetTransactionId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (PATCH /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) PatchTransactionId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID, params autogen.PatchTransactionIdParams) error {
	// TODO: implement
	return nil
}

// (PATCH /accounts/{account_id}/transactions/{transaction_id}/{item_id})
func (s *Server) PatchTransactionItemId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID, itemId autogen.UUID, params autogen.PatchTransactionItemIdParams) error {
	// TODO: implement
	return nil
}

// (GET /transactions)
func (s *Server) GetTransactions(c echo.Context, params autogen.GetTransactionsParams) error {
	// TODO: implement
	return nil
}

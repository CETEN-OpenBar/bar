package api

import (
	"bar/autogen"
	"context"
)

// (GET /transactions)
func (s *Server) GetTransactions(ctx context.Context, request autogen.GetTransactionsRequestObject) (autogen.GetTransactionsResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /account/transactions)
func (s *Server) GetCurrentAccountTransactions(ctx context.Context, request autogen.GetCurrentAccountTransactionsRequestObject) (autogen.GetCurrentAccountTransactionsResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (POST /account/transactions)
func (s *Server) PutTransactions(ctx context.Context, request autogen.PutTransactionsRequestObject) (autogen.PutTransactionsResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /accounts/{account_id}/transactions)
func (s *Server) GetAccountTransactions(ctx context.Context, request autogen.GetAccountTransactionsRequestObject) (autogen.GetAccountTransactionsResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (DELETE /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) DeleteTransactionId(ctx context.Context, request autogen.DeleteTransactionIdRequestObject) (autogen.DeleteTransactionIdResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) GetTransactionId(ctx context.Context, request autogen.GetTransactionIdRequestObject) (autogen.GetTransactionIdResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (PATCH /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) PatchTransactionId(ctx context.Context, request autogen.PatchTransactionIdRequestObject) (autogen.PatchTransactionIdResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (PATCH /accounts/{account_id}/transactions/{transaction_id}/{item_id})
func (s *Server) PatchTransactionItemId(ctx context.Context, request autogen.PatchTransactionItemIdRequestObject) (autogen.PatchTransactionItemIdResponseObject, error) {
	// TODO: implement
	return nil, nil
}

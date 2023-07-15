package api

import (
	"bar/autogen"
	"context"
)

// (GET /account)
func (s *Server) GetAccount(ctx context.Context, request autogen.GetAccountRequestObject) (autogen.GetAccountResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /accounts)
func (s *Server) GetAccounts(ctx context.Context, request autogen.GetAccountsRequestObject) (autogen.GetAccountsResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (POST /accounts)
func (s *Server) PutAccounts(ctx context.Context, request autogen.PutAccountsRequestObject) (autogen.PutAccountsResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (DELETE /accounts/{account_id})
func (s *Server) DeleteAccountId(ctx context.Context, request autogen.DeleteAccountIdRequestObject) (autogen.DeleteAccountIdResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /accounts/{account_id})
func (s *Server) GetAccountId(ctx context.Context, request autogen.GetAccountIdRequestObject) (autogen.GetAccountIdResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (PATCH /accounts/{account_id})
func (s *Server) PatchAccountId(ctx context.Context, request autogen.PatchAccountIdRequestObject) (autogen.PatchAccountIdResponseObject, error) {
	// TODO: implement
	return nil, nil
}

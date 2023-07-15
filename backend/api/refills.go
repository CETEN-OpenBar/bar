package api

import (
	"bar/autogen"
	"context"
)

// (GET /refills)
func (s *Server) GetRefills(ctx context.Context, request autogen.GetRefillsRequestObject) (autogen.GetRefillsResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /account/refills)
func (s *Server) GetSelfRefills(ctx context.Context, request autogen.GetSelfRefillsRequestObject) (autogen.GetSelfRefillsResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /accounts/{account_id}/refills)
func (s *Server) GetAccountRefills(ctx context.Context, request autogen.GetAccountRefillsRequestObject) (autogen.GetAccountRefillsResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (POST /accounts/{account_id}/refills)
func (s *Server) PostRefill(ctx context.Context, request autogen.PostRefillRequestObject) (autogen.PostRefillResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (DELETE /accounts/{account_id}/refills/{refill_id})
func (s *Server) DeleteRefill(ctx context.Context, request autogen.DeleteRefillRequestObject) (autogen.DeleteRefillResponseObject, error) {
	// TODO: implement
	return nil, nil
}

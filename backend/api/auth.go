package api

import (
	"bar/autogen"
	"context"
)

// (GET /account/qr)
func (s *Server) GetAccountQR(ctx context.Context, request autogen.GetAccountQRRequestObject) (autogen.GetAccountQRResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /auth/google/begin/{qr_nonce})
func (s *Server) ConnectAccount(ctx context.Context, request autogen.ConnectAccountRequestObject) (autogen.ConnectAccountResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /auth/google/callback)
func (s *Server) Callback(ctx context.Context, request autogen.CallbackRequestObject) (autogen.CallbackResponseObject, error) {
	// TODO: implement
	return nil, nil
}

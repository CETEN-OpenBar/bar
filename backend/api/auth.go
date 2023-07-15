package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
)

// (GET /account/qr)
func (s *Server) GetAccountQR(c echo.Context) error {
	// TODO: implement
	return nil
}

// (GET /auth/google/begin/{qr_nonce})
func (s *Server) ConnectAccount(c echo.Context, qrNonce string) error {
	// TODO: implement
	return nil
}

// (GET /auth/google/callback)
func (s *Server) Callback(c echo.Context, params autogen.CallbackParams) error {
	// TODO: implement
	return nil
}

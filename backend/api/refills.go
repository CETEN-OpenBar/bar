package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
)

// (GET /refills)
func (s *Server) GetRefills(c echo.Context, params autogen.GetRefillsParams) error {
	// TODO: implement
	return nil
}

// (GET /account/refills)
func (s *Server) GetSelfRefills(c echo.Context, params autogen.GetSelfRefillsParams) error {
	// TODO: implement
	return nil
}

// (GET /accounts/{account_id}/refills)
func (s *Server) GetAccountRefills(c echo.Context, accountId autogen.UUID, params autogen.GetAccountRefillsParams) error {
	// TODO: implement
	return nil
}

// (POST /accounts/{account_id}/refills)
func (s *Server) PostRefill(c echo.Context, accountId autogen.UUID, params autogen.PostRefillParams) error {
	// TODO: implement
	return nil
}

// (DELETE /accounts/{account_id}/refills/{refill_id})
func (s *Server) DeleteRefill(c echo.Context, accountId autogen.UUID, refillId autogen.UUID) error {
	// TODO: implement
	return nil
}

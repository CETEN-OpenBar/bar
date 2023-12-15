package api

import (
	"bar/autogen"
	"bar/internal/config"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// (POST /auth/card)
func (s *Server) ConnectPassword(c echo.Context) error {
	// Check that header "X-Local-Token" is set to the local token
	if c.Request().Header.Get("X-Local-Token") != config.GetConfig().ApiConfig.LocalToken {
		return ErrorNotAuthenticated(c)
	}

	var param autogen.ConnectPasswordJSONBody
	err := c.Bind(&param)
	if err != nil {
		return Error400(c)
	}

	account, err := s.DBackend.GetAccountByCard(c.Request().Context(), param.CardId)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return Error500(c)
		}
	}

	if account.Password == nil {
		return ErrorAccNotFound(c)
	}

	if !account.VerifyPassword(param.Password) {
		return ErrorAccNotFound(c)
	}

	s.SetCookie(c, account)

	autogen.ConnectCard200JSONResponse{
		Account: &account.Account,
	}.VisitConnectCardResponse(c.Response())
	return nil
}

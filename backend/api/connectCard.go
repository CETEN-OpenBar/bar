package api

import (
	"bar/autogen"
	"bar/internal/config"
	"bar/internal/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/mongo"
)

// (POST /auth/card)
func (s *Server) ConnectCard(c echo.Context) error {
	// Check that header "X-Local-Token" is set to the local token
	if c.Request().Header.Get("X-Local-Token") != config.GetConfig().ApiConfig.LocalToken {
		return ErrorNotAuthenticated(c)
	}

	var param autogen.ConnectCardJSONBody
	err := c.Bind(&param)
	if err != nil {
		return Error400(c)
	}

	account, err := s.DBackend.GetAccountByCard(c.Request().Context(), param.CardId)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return Error500(c)
		}
		// Create default account with 1234 pin
		account = &models.Account{
			Account: autogen.Account{
				CardId:    autogen.OptionalString(param.CardId),
				Id:        uuid.New(),
				Role:      autogen.AccountStudent,
				PriceRole: autogen.AccountPriceCeten,
				State:     autogen.AccountNotOnBoarded,
			},
		}
		account.SetPin("1234")

		_, found := onBoardCache.Get(account.Account.Id.String())
		if found {
			return ErrorAccNotFound(c)
		}

		onBoardCache.Set(account.Account.Id.String(), account, cache.DefaultExpiration)
	}

	if !account.VerifyPin(param.CardPin) {
		return ErrorAccNotFound(c)
	}

	s.SetCookie(c, account)

	autogen.ConnectCard200JSONResponse{
		Account: &account.Account,
	}.VisitConnectCardResponse(c.Response())
	return nil
}

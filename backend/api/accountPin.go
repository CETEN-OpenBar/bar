package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
)

// (PATCH /account/pin)
func (s *Server) PatchAccountPin(c echo.Context) error {
	account, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	var param autogen.PatchAccountPinJSONBody
	err = c.Bind(&param)
	if err != nil {
		return Error400(c)
	}

	if param.CardId != nil && *param.CardId != "" && account.CardId != nil && *account.CardId == "" {
		// The user doesn't have a card id yet, so we can set it without checking the pin
		account.CardId = param.CardId
		account.SetPin(param.NewCardPin)
	} else {
		// sha256 both pins
		if !account.VerifyPin(param.OldCardPin) {
			return Error400(c)
		}

		account.SetPin(param.NewCardPin)
	}

	err = s.UpdateAccount(c.Request().Context(), account)
	if err != nil {
		return Error500(c)
	}

	autogen.PatchAccountPin200JSONResponse{
		Account: &account.Account,
	}.VisitPatchAccountPinResponse(c.Response())
	return nil
}

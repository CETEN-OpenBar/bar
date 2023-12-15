package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
)

// (PATCH /account/password)
func (s *Server) PatchAccountPassword(c echo.Context) error {
	account, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	var param autogen.PatchAccountPasswordJSONBody
	err = c.Bind(&param)
	if err != nil {
		return Error400(c)
	}

	if account.Password == nil {
		// The user doesn't have a card id yet, so we can set it without checking the pin
		account.SetPin(param.NewPassword)
	} else {
		// sha256 both pins
		if !account.VerifyPassword(param.OldPassword) {
			return Error400(c)
		}

		account.SetPin(param.NewPassword)
	}

	err = s.UpdateAccount(c.Request().Context(), account)
	if err != nil {
		return Error500(c)
	}

	autogen.PatchAccountPassword200JSONResponse{
		Account: &account.Account,
	}.VisitPatchAccountPasswordResponse(c.Response())
	return nil
}

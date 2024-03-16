package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
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

// (POST /accounts/{account_id}/reset_pin)
func (s *Server) ResetAccountPin(c echo.Context, accountId autogen.UUID) error {

	admin, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	account, err := s.DBackend.GetAccount(c.Request().Context(), accountId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	// Can only reset the the pin of an account below the current role
	switch admin.Role {
	case autogen.AccountSuperAdmin:
	case autogen.AccountAdmin:
		if account.Role == autogen.AccountSuperAdmin || account.Role == autogen.AccountAdmin {
			return Error400(c)
		}
	case autogen.AccountMember:
		if account.Role == autogen.AccountSuperAdmin || account.Role == autogen.AccountAdmin || account.Role == autogen.AccountMember {
			return Error400(c)
		}
	default:
		return Error400(c)
	}

	account.SetPin("1234")

	logrus.WithField("account", account.Name()).WithField("by", admin.Name()).Info("Account pin has been reset.")
	err = s.UpdateAccount(c.Request().Context(), account)
	if err != nil {
		return Error500(c)
	}

	return autogen.ResetAccountPin200Response{}.VisitResetAccountPinResponse(c.Response())

}

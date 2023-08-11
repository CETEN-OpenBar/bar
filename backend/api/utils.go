package api

import (
	"bar/autogen"
	"bar/internal/models"
	"errors"

	"github.com/labstack/echo/v4"
)

func (s *Server) SetCookie(c echo.Context, account *models.Account) {
	if account.State != autogen.AccountNotOnBoarded {
		sess := s.getUserSess(c)
		sess.Options.MaxAge = 60 * 60 * 24 * 7 // 1 week
		sess.Options.HttpOnly = true
		sess.Options.Secure = true
		sess.Values["account_id"] = account.Account.Id.String()
		sess.Save(c.Request(), c.Response())
	} else {
		sess := s.getOnboardSess(c)
		sess.Options.MaxAge = 60 * 60 * 24 * 7 // 1 week
		sess.Options.HttpOnly = true
		sess.Options.Secure = true
		sess.Values["onboard_account_id"] = account.Account.Id.String()
		sess.Save(c.Request(), c.Response())
	}

	if account.IsAdmin() {
		sess := s.getAdminSess(c)
		sess.Options.MaxAge = 60 * 60 * 24 * 7 // 1 week
		sess.Options.HttpOnly = true
		sess.Options.Secure = true
		sess.Values["admin_account_id"] = account.Account.Id.String()
		sess.Save(c.Request(), c.Response())
	}
}

func (s *Server) RemoveCookies(c echo.Context) {
	sess := s.getUserSess(c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	sess = s.getAdminSess(c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	sess = s.getOnboardSess(c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())
}

func (s *Server) RemoveOnBoardCookie(c echo.Context) {
	sess := s.getOnboardSess(c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())
}

func MustGetUser(c echo.Context) (*models.Account, error) {
	logged := c.Get("userLogged").(bool)
	if !logged {
		ErrorNotAuthenticated(c)
		return nil, errors.New("not authenticated")
	}

	account := c.Get("userAccount").(*models.Account)
	return account, nil
}

func MustGetAdmin(c echo.Context) (*models.Account, error) {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		ErrorForbidden(c)
		return nil, errors.New("not authenticated")
	}

	account := c.Get("adminAccount").(*models.Account)
	return account, nil
}

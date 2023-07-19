package api

import (
	"bar/autogen"
	"bar/internal/models"

	"github.com/labstack/echo/v4"
)

func Error500(c echo.Context) error {
	resp := autogen.GetAccounts500JSONResponse{
		Message:   autogen.MsgInternalServerError,
		ErrorCode: autogen.ErrInternalServerError,
	}
	resp.VisitGetAccountsResponse(c.Response())
	return nil
}

func Error400(c echo.Context) error {
	resp := autogen.ConnectAccount400JSONResponse{
		Message:   autogen.MsgBadRequest,
		ErrorCode: autogen.ErrBadRequest,
	}
	resp.VisitConnectAccountResponse(c.Response())
	return nil
}

func ErrorNotAuthenticated(c echo.Context) error {
	resp := autogen.GetAccounts401JSONResponse{
		Message:   autogen.MsgNotAuthenticated,
		ErrorCode: autogen.ErrNotAuthenticated,
	}
	resp.VisitGetAccountsResponse(c.Response())
	return nil
}

func Error409(c echo.Context) error {
	resp := autogen.PostAccounts409JSONResponse{
		Message:   autogen.MsgAccountAlreadyExists,
		ErrorCode: autogen.ErrBadRequest,
	}
	resp.VisitPostAccountsResponse(c.Response())
	return nil
}

func ErrorAccNotFound(c echo.Context) error {
	resp := autogen.GetAccount401JSONResponse{
		Message:   autogen.MsgAccountNotFound,
		ErrorCode: autogen.ErrAccountNotFound,
	}
	resp.VisitGetAccountResponse(c.Response())
	return nil
}

func (s *Server) SetCookie(c echo.Context, account *models.Account) {
	sess := s.getUserSess(c)
	sess.Options.MaxAge = 60 * 60 * 24 * 7 // 1 week
	sess.Options.HttpOnly = true
	sess.Options.Secure = true
	sess.Values["account_id"] = account.Account.Id.String()
	sess.Save(c.Request(), c.Response())

	if account.IsAdmin() {
		sess := s.getAdminSess(c)
		sess.Options.MaxAge = 60 * 60 * 24 * 7 // 1 week
		sess.Options.HttpOnly = true
		sess.Options.Secure = true
		sess.Values["admin_account_id"] = account.Account.Id.String()
		if account.Role == autogen.AccountSuperAdmin {
			sess.Values["super_admin"] = true
		}
		sess.Save(c.Request(), c.Response())
	}
}

func (s *Server) RemoveCookie(c echo.Context) {
	sess := s.getUserSess(c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	sess = s.getAdminSess(c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())
}

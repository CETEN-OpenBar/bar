package api

import (
	"bar/autogen"

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

func Error401(c echo.Context) error {
	resp := autogen.GetAccounts401JSONResponse{
		Message:   autogen.MsgAccountNotAllowed,
		ErrorCode: autogen.ErrForbidden,
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

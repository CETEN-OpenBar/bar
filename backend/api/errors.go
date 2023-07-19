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

func ErrorImageNotFound(c echo.Context) error {
	resp := autogen.GetCarouselImage404JSONResponse{
		Message:   autogen.MsgImageNotFound,
		ErrorCode: autogen.ErrImageNotFound,
	}
	resp.VisitGetCarouselImageResponse(c.Response())
	return nil
}

func ErrorTextNotFound(c echo.Context) error {
	resp := autogen.MarkDeleteCarouselText404JSONResponse{
		Message:   autogen.MsgTextNotFound,
		ErrorCode: autogen.ErrTextNotFound,
	}
	resp.VisitMarkDeleteCarouselTextResponse(c.Response())
	return nil
}

func ErrorCategoryNotFound(c echo.Context) error {
	resp := autogen.GetCategory404JSONResponse{
		Message:   autogen.MsgCategoryNotFound,
		ErrorCode: autogen.ErrCategoryNotFound,
	}
	resp.VisitGetCategoryResponse(c.Response())
	return nil
}

func ErrorItemNotFound(c echo.Context) error {
	resp := autogen.GetCategoryItems404JSONResponse{
		Message:   autogen.MsgItemNotFound,
		ErrorCode: autogen.ErrItemNotFound,
	}
	resp.VisitGetCategoryItemsResponse(c.Response())
	return nil
}

func ErrorRefillNotFound(c echo.Context) error {
	resp := autogen.DeleteRefill404JSONResponse{
		Message:   autogen.MsgRefillNotFound,
		ErrorCode: autogen.ErrRefillNotFound,
	}
	resp.VisitDeleteRefillResponse(c.Response())
	return nil
}

func ErrorTransactionNotFound(c echo.Context) error {
	resp := autogen.DeleteTransaction404JSONResponse{
		Message:   autogen.MsgRefillNotFound,
		ErrorCode: autogen.ErrRefillNotFound,
	}
	resp.VisitDeleteTransactionResponse(c.Response())
	return nil
}

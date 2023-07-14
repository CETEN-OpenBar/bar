package api

import (
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func (s *Server) AuthMiddleware(f runtime.StrictEchoHandlerFunc, operationID string) runtime.StrictEchoHandlerFunc {
	return func(c echo.Context, req interface{}) (resp interface{}, err error) {
		s := c.Get("userStore")
		userStore, ok := s.(sessions.Store)
		if !ok {
			return nil, echo.NewHTTPError(500, "userStore not found")
		}
		s = c.Get("adminStore")
		adminStore, ok := s.(sessions.Store)
		if !ok {
			return nil, echo.NewHTTPError(500, "adminStore not found")
		}

		userSess, err := userStore.Get(c.Request(), "BAR_SESS")
		if err != nil {
			return nil, echo.NewHTTPError(500, "session not found")
		}
		adminSess, err := adminStore.Get(c.Request(), "BAR_ADMIN_SESS")
		if err != nil {
			return nil, echo.NewHTTPError(500, "session not found")
		}

		c.Set("userSess", userSess)
		c.Set("adminSess", adminSess)

		return f(c, req)
	}
}

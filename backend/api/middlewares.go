package api

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		v := c.Get("userStore")
		userStore, ok := v.(sessions.Store)
		if !ok {
			return echo.NewHTTPError(500, "userStore not found")
		}
		v = c.Get("adminStore")
		adminStore, ok := v.(sessions.Store)
		if !ok {
			return echo.NewHTTPError(500, "adminStore not found")
		}

		userSess, err := userStore.Get(c.Request(), "BAR_SESS")
		if err != nil {
			return echo.NewHTTPError(500, "session not found")
		}
		adminSess, err := adminStore.Get(c.Request(), "BAR_ADMIN_SESS")
		if err != nil {
			return echo.NewHTTPError(500, "session not found")
		}

		c.Set("userSess", userSess)
		c.Set("adminSess", adminSess)

		c.Set("userLogged", false)
		c.Set("adminLogged", false)

		// Get user account from cookie
		accountID, ok := userSess.Values["account_id"].(string)
		if ok {
			// Get account from database
			account, err := s.DBackend.GetAccount(c.Request().Context(), accountID)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					// Delete cookie
					userSess.Options.MaxAge = -1
					userSess.Save(c.Request(), c.Response())
					return ErrorAccNotFound(c)
				}
				logrus.Error(err)
				return Error500(c)
			}

			c.Set("userLogged", true)
			c.Set("userAccountID", accountID)
			c.Set("userAccount", account)
		}

		// Get admin account from cookie
		adminId, ok := adminSess.Values["admin_account_id"].(string)
		if ok {
			// Get account from database
			account, err := s.DBackend.GetAccount(c.Request().Context(), adminId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					// Delete cookie
					adminSess.Options.MaxAge = -1
					adminSess.Save(c.Request(), c.Response())
					return ErrorAccNotFound(c)
				}
				logrus.Error(err)
				return Error500(c)
			}

			c.Set("adminLogged", true)
			c.Set("adminAccountID", adminId)
			c.Set("adminAccount", account)
			c.Set("adminAccountRole", account.Role)
		}

		return next(c)
	}
}

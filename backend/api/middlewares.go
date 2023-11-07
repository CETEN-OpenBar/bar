package api

import (
	"bar/internal/models"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var onBoardCache = cache.New(5*time.Minute, 10*time.Minute)

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
		v = c.Get("onBoardStore")
		onBoardStore, ok := v.(sessions.Store)
		if !ok {
			return echo.NewHTTPError(500, "onBoardStore not found")
		}

		userSess, err := userStore.Get(c.Request(), "BAR_SESS")
		if err != nil {
			return echo.NewHTTPError(500, "session not found")
		}
		adminSess, err := adminStore.Get(c.Request(), "BAR_ADMIN_SESS")
		if err != nil {
			return echo.NewHTTPError(500, "session not found")
		}
		onBoardSess, err := onBoardStore.Get(c.Request(), "BAR_ONBOARD_SESS")
		if err != nil {
			return echo.NewHTTPError(500, "session not found")
		}

		c.Set("userSess", userSess)
		c.Set("adminSess", adminSess)
		c.Set("onBoardSess", onBoardSess)

		c.Set("userLogged", false)
		c.Set("adminLogged", false)
		c.Set("onBoardLogged", false)

		onBoardID, ok := onBoardSess.Values["onboard_account_id"].(string)
		if ok {
			// Get account from database
			acc, found := onBoardCache.Get(onBoardID)
			if !found {
				// Remove cookie and go on
				onBoardSess.Options.MaxAge = -1
				onBoardSess.Save(c.Request(), c.Response())
			} else {
				account := acc.(*models.Account)
				c.Set("onBoardLogged", true)
				c.Set("onBoardAccountID", onBoardID)
				c.Set("onBoardAccount", account)
			}
		}

		// Get user account from cookie
		accountID, ok := userSess.Values["account_id"].(string)
		if ok {
			// Get account from database
			account, err := s.DBackend.GetAccount(c.Request().Context(), accountID)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					// Delete cookie
					userSess.Options.MaxAge = -1
					adminSess.Options.MaxAge = -1
					userSess.Save(c.Request(), c.Response())
					adminSess.Save(c.Request(), c.Response())
					return ErrorAccNotFound(c)
				}
				logrus.Error(err)
				return Error500(c)
			}

			if account.IsBlocked() {
				logrus.Warnf("Account %s is blocked", accountID)
				// Delete cookie
				userSess.Options.MaxAge = -1
				adminSess.Options.MaxAge = -1
				userSess.Save(c.Request(), c.Response())
				adminSess.Save(c.Request(), c.Response())
				return Error403(c)
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

			if account.IsAdmin() {
				c.Set("adminLogged", true)
				c.Set("adminAccountID", adminId)
				c.Set("adminAccount", account)
				c.Set("adminAccountRole", account.Role)
			}
		}

		return next(c)
	}
}

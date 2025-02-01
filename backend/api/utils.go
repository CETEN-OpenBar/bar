package api

import (
	"bar/autogen"
	"bar/internal/models"
	"errors"
	"math"

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

func MustGetUserOrOnBoard(c echo.Context) (*models.Account, error) {
	logged := c.Get("userLogged").(bool)
	loggedOnBoard := c.Get("onBoardLogged").(bool)
	if !logged && !loggedOnBoard {
		ErrorNotAuthenticated(c)
		return nil, errors.New("not authenticated")
	}

	var account *models.Account

	if logged {
		account = c.Get("userAccount").(*models.Account)
	}

	if loggedOnBoard {
		account = c.Get("onBoardAccount").(*models.Account)
	}

	return account, nil
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

	if account.State == autogen.AccountNotOnBoarded {
		ErrorForbidden(c)
		return nil, errors.New("not authenticated")
	}

	return account, nil
}

func UpdateItem(item *models.Item, category *models.Category, restockItem autogen.RestockItem) *models.Item {
	item.State = autogen.ItemBuyable
	item.AmountLeft += restockItem.AmountOfBundle * restockItem.AmountPerBundle
	item.LastTva = &restockItem.Tva
	if !category.SpecialPrice {
		item.Prices.Coutant = uint64(math.Ceil(float64(restockItem.BundleCostTtc) / (float64(restockItem.AmountPerBundle))))
		if item.Prices.Coutant < 30 {
			item.Prices.Externe = arrondiAuMutilple(item.Prices.Coutant, 5) + 20
			item.Prices.Ceten = arrondiAuMutilple(item.Prices.Coutant, 5) + 10
			item.Prices.StaffBar = arrondiAuMutilple(item.Prices.Coutant, 5) + 5
			item.Prices.Privilegies = arrondiAuMutilple(item.Prices.Coutant, 5) + 5
			item.Prices.Menu = arrondiAuMutilple(item.Prices.Coutant, 5) + 10
		} else if item.Prices.Coutant >= 30 && item.Prices.Coutant < 130 {
			item.Prices.Externe = arrondiAuMutilple(item.Prices.Coutant*3/2, 5)
			item.Prices.Ceten = arrondiAuMutilple(item.Prices.Coutant*113/100, 5)
			item.Prices.StaffBar = arrondiAuMutilple(item.Prices.Coutant*108/100, 5)
			item.Prices.Privilegies = arrondiAuMutilple(item.Prices.Coutant*11/10, 5)
			item.Prices.Menu = arrondiAuMutilple(item.Prices.Coutant*13/10, 5)
		} else if item.Prices.Coutant >= 130 && item.Prices.Coutant < 300 {
			item.Prices.Externe = arrondiAuMutilple(item.Prices.Coutant*14/10, 5)
			item.Prices.Ceten = arrondiAuMutilple(item.Prices.Coutant*11/10, 5)
			item.Prices.StaffBar = arrondiAuMutilple(item.Prices.Coutant*108/100, 5)
			item.Prices.Privilegies = arrondiAuMutilple(item.Prices.Coutant*11/10, 5)
			item.Prices.Menu = arrondiAuMutilple(item.Prices.Coutant*12/10, 5)
		} else if item.Prices.Coutant >= 300 {
			item.Prices.Externe = arrondiAuMutilple(item.Prices.Coutant*125/100, 5)
			item.Prices.Ceten = arrondiAuMutilple(item.Prices.Coutant*11/10, 5)
			item.Prices.StaffBar = arrondiAuMutilple(item.Prices.Coutant*105/100, 5)
			item.Prices.Privilegies = arrondiAuMutilple(item.Prices.Coutant*11/10, 5)
			item.Prices.Menu = arrondiAuMutilple(item.Prices.Coutant*1125/1000, 5)
		}
	} else {
		item.Prices.Coutant = uint64(math.Ceil(float64(restockItem.BundleCostTtc) / (float64(restockItem.AmountPerBundle))))
		if item.Prices.Coutant < 30 {
			item.Prices.Externe = arrondiAuMutilple(item.Prices.Coutant, 5) + 20
			item.Prices.Ceten = arrondiAuMutilple(item.Prices.Coutant, 5) + 10
			item.Prices.StaffBar = item.Prices.Ceten
			item.Prices.Privilegies = item.Prices.Ceten
			item.Prices.Menu = item.Prices.Ceten
		} else if item.Prices.Coutant >= 30 && item.Prices.Coutant < 130 {
			item.Prices.Externe = arrondiAuMutilple(item.Prices.Coutant*3/2, 5)
			item.Prices.Ceten = arrondiAuMutilple(item.Prices.Coutant*113/100, 5)
			item.Prices.StaffBar = item.Prices.Ceten
			item.Prices.Privilegies = item.Prices.Ceten
			item.Prices.Menu = item.Prices.Ceten
		} else if item.Prices.Coutant >= 130 && item.Prices.Coutant < 300 {
			item.Prices.Externe = arrondiAuMutilple(item.Prices.Coutant*14/10, 5)
			item.Prices.Ceten = arrondiAuMutilple(item.Prices.Coutant*11/10, 5)
			item.Prices.StaffBar = item.Prices.Ceten
			item.Prices.Privilegies = item.Prices.Ceten
			item.Prices.Menu = item.Prices.Ceten
		} else if item.Prices.Coutant >= 300 {
			item.Prices.Externe = arrondiAuMutilple(item.Prices.Coutant*125/100, 5)
			item.Prices.Ceten = arrondiAuMutilple(item.Prices.Coutant*11/10, 5)
			item.Prices.StaffBar = item.Prices.Ceten
			item.Prices.Privilegies = item.Prices.Ceten
			item.Prices.Menu = item.Prices.Ceten
		}
		item.Prices.Coutant = item.Prices.Ceten
	}
	return item
}
package api

import (
	"bar/autogen"
	"bar/internal/models"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (POST /account/transactions)
func (s *Server) PostTransactions(c echo.Context) error {
	// Get account from cookie
	_, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	accountID := c.Get("userAccountID").(string)
	account := c.Get("userAccount").(*models.Account)

	transaction := &models.Transaction{
		Transaction: autogen.Transaction{
			AccountId:   accountID,
			AccountName: account.Name(),
			State:       autogen.TransactionStarted,
			Id:          uuid.New(),
		},
	}

	var potentialTransaction autogen.NewTransaction

	// Check that pin matches
	err = c.Bind(&potentialTransaction)
	if err != nil {
		logrus.Error(err)
		return Error400(c)
	}

	if potentialTransaction.Items == nil || len(potentialTransaction.Items) == 0 {
		logrus.Warnf("Transaction %s has no items", transaction.Id.String())
		return Error400(c)
	}

	if !account.VerifyPin(potentialTransaction.CardPin) {
		autogen.PostTransactions401JSONResponse{
			Message:   autogen.MsgNotAuthenticated,
			ErrorCode: autogen.ErrNotAuthenticated,
		}.VisitPostTransactionsResponse(c.Response())
		return nil
	}

	_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		var transactionCost uint64
		var fetchedItems = make(map[string]*models.Item)

		// Fill up fetchedItems first to avoid multiple updates
		for _, potentialItem := range potentialTransaction.Items {
			item, ok := fetchedItems[potentialItem.ItemId.String()]
			if !ok {
				// Verify that item exists, can be bought, is in stock, and can be bought for that amount
				item, err = s.DBackend.GetItem(c.Request().Context(), potentialItem.ItemId.String())
				if err != nil {
					if err == mongo.ErrNoDocuments {
						return nil, ErrorItemNotFound(c)
					}
					logrus.Error(err)
					return nil, Error500(c)
				}
				fetchedItems[potentialItem.ItemId.String()] = item
			} else {
				logrus.Warnf("Item %s already fetched", item.Id.String())
				return nil, Error400(c)
			}
		}

		for _, potentialItem := range potentialTransaction.Items {
			item := fetchedItems[potentialItem.ItemId.String()]
			if !item.IsMenu {
				continue
			}
			for _, menuItem := range *item.MenuItems {
				mItem, err := s.DBackend.GetItem(c.Request().Context(), menuItem.Id.String())
				if err != nil {
					if err == mongo.ErrNoDocuments {
						return nil, ErrorItemNotFound(c)
					}
					logrus.Error(err)
					return nil, Error500(c)
				}
				fetchedItems[mItem.Id.String()] = mItem
			}
			if potentialItem.PickedCategoriesItems == nil && item.MenuCategories != nil && len(*item.MenuCategories) > 0 {
				logrus.Warnf("Menu item %s has no picked categories", item.Id.String())
				return nil, Error400(c)
			}
			for _, pickedItem := range *potentialItem.PickedCategoriesItems {
				mItem, err := s.DBackend.GetItem(c.Request().Context(), pickedItem.ItemId.String())
				if err != nil {
					if err == mongo.ErrNoDocuments {
						return nil, ErrorItemNotFound(c)
					}
					logrus.Error(err)
					return nil, Error500(c)
				}
				fetchedItems[mItem.Id.String()] = mItem
			}
		}
		// Finished filling up fetchedItems

		// Check for menus in transaction first, add thos item without a price to the transaction
		for _, potentialItem := range potentialTransaction.Items {
			item := fetchedItems[potentialItem.ItemId.String()]

			if !item.IsMenu {
				continue
			}

			t := uint64(time.Since(time.Now().Truncate(24 * time.Hour)).Seconds())

			if item.AvailableFrom != nil && *item.AvailableFrom > t && item.AvailableUntil != nil && *item.AvailableUntil < t {
				logrus.Warnf("Menu item %s is not available", item.Id.String())
				return nil, Error400(c)
			}

			for _, menuItem := range *item.MenuItems {
				mItem, err := s.DBackend.GetItem(c.Request().Context(), menuItem.Id.String())
				if err != nil {
					if err == mongo.ErrNoDocuments {
						return nil, ErrorItemNotFound(c)
					}
					logrus.Error(err)
					return nil, Error500(c)
				}
				fetchedItems[mItem.Id.String()] = mItem

				if mItem.State == autogen.ItemNotBuyable {
					logrus.Warnf("Menu item %s is not buyable", item.Id.String())
					return nil, Error400(c)
				}
				if mItem.AmountLeft < menuItem.Amount {
					logrus.Warnf("Menu item %s is not in stock", item.Id.String())
					return nil, Error400(c)
				}
				if mItem.BuyLimit != nil {
					if *item.BuyLimit < menuItem.Amount {
						logrus.Warnf("Menu item %s cannot be bought for that amount", item.Id.String())
						return nil, Error400(c)
					}
				}
				mItem.AmountLeft -= menuItem.Amount
				menuItem.Name = mItem.Name
				menuItem.PictureUri = mItem.PictureUri
			}

			// Check that there's no item that would be in other categories
			for _, pickedItem := range *potentialItem.PickedCategoriesItems {
				pItem := fetchedItems[pickedItem.ItemId.String()]
				found := false
				for _, menuCategory := range *item.MenuCategories {
					if pItem.CategoryId.String() == menuCategory.Id.String() {
						found = true
						break
					}
				}
				if !found {
					logrus.Warnf("Item %s is not in menu but was supplied", pItem.Id.String())
					return nil, Error400(c)
				}
			}

			// Check that all categories are not picked more than the amount
			for _, menuCategory := range *item.MenuCategories {
				amountPicked := uint64(0)
				for _, pickedItem := range *potentialItem.PickedCategoriesItems {
					pItem := fetchedItems[pickedItem.ItemId.String()]
					if pItem.CategoryId.String() == menuCategory.Id.String() {
						amountPicked += pickedItem.Amount
						pItem.AmountLeft -= pickedItem.Amount
					}
				}
				if amountPicked > menuCategory.Amount {
					logrus.Warnf("Menu category %s cannot be picked for that amount", menuCategory.Id.String())
					return nil, Error400(c)
				}
			}
		}

		for _, potentialItem := range potentialTransaction.Items {
			item := fetchedItems[potentialItem.ItemId.String()]

			if item.State == autogen.ItemNotBuyable {
				logrus.Warnf("Item %s is not buyable", item.Id.String())
				return nil, Error400(c)
			}
			if item.AmountLeft < potentialItem.Amount {
				logrus.Warnf("Item %s is not in stock", item.Id.String())
				return nil, Error400(c)
			}
			if item.BuyLimit != nil {
				if *item.BuyLimit < potentialItem.Amount {
					logrus.Warnf("Item %s cannot be bought for that amount", item.Id.String())
					return nil, Error400(c)
				}
			}

			t := uint64(time.Since(time.Now().Truncate(24 * time.Hour)).Seconds())

			if item.AvailableFrom != nil && *item.AvailableFrom > t && item.AvailableUntil != nil && *item.AvailableUntil < t {
				logrus.Warnf("Menu item %s is not available", item.Id.String())
				return nil, Error400(c)
			}

			var picked []autogen.TransactionItem
			if potentialItem.PickedCategoriesItems != nil && item.IsMenu {
				for _, pickedItem := range *potentialItem.PickedCategoriesItems {
					i := fetchedItems[pickedItem.ItemId.String()]
					picked = append(picked, autogen.TransactionItem{
						ItemAmount:     pickedItem.Amount,
						ItemId:         pickedItem.ItemId,
						ItemName:       i.Name,
						PictureUri:     i.PictureUri,
						State:          autogen.TransactionItemStarted,
						IsMenu:         i.IsMenu,
						MenuItems:      i.MenuItems,
						MenuCategories: i.MenuCategories,
					})
				}
			}

			transaction.Items = append(transaction.Items, autogen.TransactionItem{
				ItemAmount:            potentialItem.Amount,
				ItemId:                potentialItem.ItemId,
				ItemName:              item.Name,
				PictureUri:            item.PictureUri,
				State:                 autogen.TransactionItemStarted,
				IsMenu:                item.IsMenu,
				MenuItems:             item.MenuItems,
				MenuCategories:        item.MenuCategories,
				PickedCategoriesItems: &picked,
				UnitCost:              item.RealPrice(account.PriceRole),
				TotalCost:             item.RealPrice(account.PriceRole) * potentialItem.Amount,
			})

			transactionCost += item.RealPrice(account.PriceRole) * potentialItem.Amount
			item.AmountLeft -= potentialItem.Amount
		}

		transaction.TotalCost = transactionCost

		// update account balance
		if int64(transactionCost) > account.Points {
			account.Balance -= int64(transactionCost) - account.Points
			account.Points = 0
		} else {
			account.Points -= int64(transactionCost)
		}

		if account.Role != autogen.AccountGhost {
			if account.Balance < 0 {
				logrus.Warnf("Account %s does not have enough money", accountID)
				return nil, Error400(c)
			}
		}

		err = s.DBackend.CreateTransaction(ctx, transaction)
		if err != nil {
			logrus.Error(err)
			return nil, errors.New("failed to create transaction")
		}

		err = s.DBackend.UpdateAccount(ctx, account)
		if err != nil {
			logrus.Error(err)
			return nil, errors.New("failed to update account")
		}

		// update items
		for _, item := range fetchedItems {
			err = s.DBackend.UpdateItem(ctx, item)
			if err != nil {
				logrus.Error(err)
				return nil, errors.New("failed to update item")
			}
		}

		return nil, nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}

	if c.Response().Committed {
		return nil
	}

	autogen.PostTransactions201JSONResponse(transaction.Transaction).VisitPostTransactionsResponse(c.Response())
	return nil
}

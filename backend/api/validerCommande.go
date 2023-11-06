package api

import (
	"bar/autogen"
	"bar/internal/models"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (PATCH /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) PatchTransactionId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID, params autogen.PatchTransactionIdParams) error {
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	account := c.Get("userAccount").(*models.Account)

	// Get transaction from database
	transaction, err := s.DBackend.GetTransaction(c.Request().Context(), transactionId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorTransactionNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	oldState := transaction.State

	if oldState != autogen.TransactionCanceled && params.State == autogen.TransactionCanceled {
		transaction.State = params.State
		_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {

			// update account balance
			account.Balance += int64(transaction.TotalCost)
			err = s.DBackend.UpdateAccount(ctx, account)
			if err != nil {
				logrus.Error(err)
				return nil, errors.New("failed to update account")
			}

			// update items
			for i, txitem := range transaction.Items {
				if txitem.State == autogen.TransactionItemCanceled {
					continue
				}

				item, err := s.DBackend.GetItem(ctx, txitem.ItemId.String())
				if err != nil {
					continue
				}

				item.AmountLeft += txitem.ItemAmount

				if txitem.IsMenu {
					for _, subitem := range *txitem.Items {
						sItem, err := s.DBackend.GetItem(ctx, subitem.Id.String())
						if err != nil {
							continue
						}
						sItem.AmountLeft += subitem.Amount
						err = s.DBackend.UpdateItem(ctx, sItem)
						if err != nil {
							logrus.Error(err)
							return nil, errors.New("failed to update item")
						}
					}
				}

				err = s.DBackend.UpdateItem(ctx, item)
				if err != nil {
					logrus.Error(err)
					return nil, errors.New("failed to update item")
				}

				txitem.State = autogen.TransactionItemCanceled
				txitem.TotalCost = 0
				transaction.Items[i] = txitem
			}

			transaction.TotalCost = 0

			err = s.DBackend.UpdateTransaction(ctx, transaction)
			if err != nil {
				logrus.Error(err)
				return nil, errors.New("failed to create transaction")
			}

			return nil, nil
		})
		if err != nil {
			logrus.Error(err)
			return Error500(c)
		}
	} else if oldState == autogen.TransactionCanceled && params.State != autogen.TransactionCanceled {
		logrus.Error("Cannot validate a canceled transaction")
		return Error400(c)
	} else {
		transaction.State = params.State
		err = s.DBackend.UpdateTransaction(c.Request().Context(), transaction)
		if err != nil {
			logrus.Error(err)
			return Error500(c)
		}
	}

	return nil
}

// (PATCH /accounts/{account_id}/transactions/{transaction_id}/{item_id})
func (s *Server) PatchTransactionItemId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID, itemId autogen.UUID, params autogen.PatchTransactionItemIdParams) error {
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	account := c.Get("userAccount").(*models.Account)

	// Get transaction from database
	transaction, err := s.DBackend.GetTransaction(c.Request().Context(), transactionId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorTransactionNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	if transaction.State == autogen.TransactionFinished {
		return Error400(c)
	}

	var item *autogen.TransactionItem

	for i, titem := range transaction.Items {
		if titem.ItemId == itemId {
			item = &transaction.Items[i]
			break
		}
	}

	oldState := item.State
	oldAmount := item.ItemAmount
	oldCost := item.TotalCost

	if params.State != nil {
		item.State = *params.State
	} else if params.Amount != nil {
		if *params.Amount > oldAmount {
			return Error400(c)
		}
		item.ItemAmount = *params.Amount
		item.TotalCost = *params.Amount * item.UnitCost
	}

	if params.AlreadyDone != nil {
		if *params.AlreadyDone > item.ItemAmount {
			return Error400(c)
		}
		item.ItemAlreadyDone = *params.AlreadyDone
		if item.ItemAlreadyDone == item.ItemAmount {
			item.State = autogen.TransactionItemFinished
		}
	}

	origItem, err := s.DBackend.GetItem(c.Request().Context(), itemId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorItemNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	if oldState == autogen.TransactionItemCanceled && item.State != autogen.TransactionItemCanceled {
		// Can't do that
		return Error400(c)
	}

	_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		if oldState != autogen.TransactionItemCanceled && item.State == autogen.TransactionItemCanceled {
			origItem.AmountLeft += item.ItemAmount
			account.Balance += int64(item.TotalCost)
			transaction.TotalCost -= item.TotalCost
			item.TotalCost = 0

			if item.IsMenu {
				for _, subitem := range *item.Items {
					origSubItem, err := s.DBackend.GetItem(c.Request().Context(), subitem.Id.String())
					if err != nil {
						if err == mongo.ErrNoDocuments {
							continue
						}
						logrus.Error(err)
						return nil, errors.New("failed to get item")
					}
					origSubItem.AmountLeft += subitem.Amount
					err = s.DBackend.UpdateItem(c.Request().Context(), origSubItem)
					if err != nil {
						logrus.Error(err)
						return nil, errors.New("failed to update item")
					}
				}
			}
		} else {
			origItem.AmountLeft += oldAmount - item.ItemAmount
			account.Balance += int64(oldCost - item.TotalCost)
			transaction.TotalCost += item.TotalCost - oldCost

			if item.IsMenu {
				for _, subitem := range *item.Items {
					origSubItem, err := s.DBackend.GetItem(c.Request().Context(), subitem.Id.String())
					if err != nil {
						if err == mongo.ErrNoDocuments {
							continue
						}
						logrus.Error(err)
						return nil, errors.New("failed to get item")
					}
					origSubItem.AmountLeft += subitem.Amount
					err = s.DBackend.UpdateItem(c.Request().Context(), origSubItem)
					if err != nil {
						logrus.Error(err)
						return nil, errors.New("failed to update item")
					}
				}
			}
		}

		err = s.DBackend.UpdateTransaction(ctx, transaction)
		if err != nil {
			logrus.Error(err)
			return nil, errors.New("failed to update transaction")
		}

		err = s.DBackend.UpdateAccount(ctx, account)
		if err != nil {
			logrus.Error(err)
			return nil, errors.New("failed to update account")
		}

		err = s.DBackend.UpdateItem(ctx, origItem)
		if err != nil {
			logrus.Error(err)
			return nil, errors.New("failed to update item")
		}
		return nil, nil
	})

	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}
	return nil
}

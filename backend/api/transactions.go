package api

import (
	"bar/autogen"
	"bar/internal/models"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (POST /account/transactions)
func (s *Server) PostTransactions(c echo.Context) error {
	// Get account from cookie
	logged := c.Get("userLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	accountID := c.Get("userAccountID").(string)
	account := c.Get("userAccount").(*models.Account)

	transaction := &models.Transaction{
		Transaction: autogen.Transaction{
			AccountId: accountID,
			State:     autogen.TransactionStarted,
			Id:        uuid.New(),
		},
	}

	var potentialTransaction autogen.NewTransaction
	var fetchedItems = make(map[string]*models.Item)

	// Check that pin matches
	err := c.Bind(&potentialTransaction)
	if err != nil {
		logrus.Error(err)
		return Error400(c)
	}

	if fmt.Sprintf("%x", sha256.Sum256([]byte(potentialTransaction.CardPin))) != account.CardPin {
		autogen.PostTransactions401JSONResponse{
			Message:   autogen.MsgNotAuthenticated,
			ErrorCode: autogen.ErrNotAuthenticated,
		}.VisitPostTransactionsResponse(c.Response())
		return nil
	}

	var transactionCost uint64
	for _, potentialItem := range potentialTransaction.Items {
		item, ok := fetchedItems[potentialItem.ItemId.String()]
		if !ok {
			// Verify that item exists, can be bought, is in stock, and can be bought for that amount
			item, err = s.DBackend.GetItem(c.Request().Context(), potentialItem.ItemId.String())
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return ErrorItemNotFound(c)
				}
				logrus.Error(err)
				return Error500(c)
			}
			fetchedItems[potentialItem.ItemId.String()] = item
		}

		if item.State == autogen.ItemNotBuyable {
			return Error400(c)
		}
		if item.AmountLeft < potentialItem.Amount {
			return Error400(c)
		}
		item.AmountLeft -= potentialItem.Amount
		if item.BuyLimit < potentialItem.Amount {
			return Error400(c)
		}

		transaction.Items = append(transaction.Items, autogen.TransactionItem{
			ItemAmount: potentialItem.Amount,
			ItemId:     potentialItem.ItemId,
			PictureUri: item.PictureUri,
			State:      autogen.TransactionItemStarted,
			UnitCost:   item.Price,
			TotalCost:  item.Price * potentialItem.Amount,
		})

		transactionCost += item.Price * potentialItem.Amount
	}

	transaction.TotalCost = transactionCost

	_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
		err = s.DBackend.CreateTransaction(ctx, transaction)
		if err != nil {
			logrus.Error(err)
			return nil, errors.New("failed to create transaction")
		}

		// update account balance
		account.Balance -= int64(transactionCost)
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
		return Error500(c)
	}

	autogen.PostTransactions201JSONResponse(transaction.Transaction).VisitPostTransactionsResponse(c.Response())
	return nil
}

// (GET /accounts/{account_id}/transactions)
func (s *Server) GetAccountTransactions(c echo.Context, accountId autogen.UUID, params autogen.GetAccountTransactionsParams) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	var page uint64
	if params.Page != nil {
		page = *params.Page
	}
	if page > 0 {
		page--
	}

	var limit uint64 = 10
	if params.Limit != nil {
		limit = *params.Limit
	}
	if limit > 100 {
		limit = 100
	}

	var state string
	if params.State != nil {
		state = string(*params.State)
	}

	count, err := s.DBackend.CountTransactions(c.Request().Context(), accountId.String(), state)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	maxPage := count / limit

	if page > maxPage {
		page = maxPage
	}

	data, err := s.DBackend.GetTransactions(c.Request().Context(), accountId.String(), page, limit, state)
	if err != nil {
		return Error500(c)
	}

	transactions := make([]autogen.Transaction, len(data))
	for i, transaction := range data {
		transactions[i] = transaction.Transaction
	}

	autogen.GetAccountTransactions200JSONResponse{
		Transactions: &transactions,
		Limit:        &limit,
		Page:         &page,
		MaxPage:      &maxPage,
	}.VisitGetAccountTransactionsResponse(c.Response())
	return nil
}

// (GET /account/transactions)
func (s *Server) GetCurrentAccountTransactions(c echo.Context, params autogen.GetCurrentAccountTransactionsParams) error {
	// Get account from cookie
	logged := c.Get("userLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	accountID := c.Get("userAccountID").(string)

	var page uint64
	if params.Page != nil {
		page = *params.Page
	}
	if page > 0 {
		page--
	}

	var limit uint64 = 10
	if params.Limit != nil {
		limit = *params.Limit
	}
	if limit > 100 {
		limit = 100
	}

	var state string
	if params.State != nil {
		state = string(*params.State)
	}

	count, err := s.DBackend.CountTransactions(c.Request().Context(), accountID, state)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	maxPage := count / limit

	if page > maxPage {
		page = maxPage
	}

	data, err := s.DBackend.GetTransactions(c.Request().Context(), accountID, page, limit, state)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	transactions := make([]autogen.Transaction, len(data))
	for i, transaction := range data {
		transactions[i] = transaction.Transaction
	}

	autogen.GetCurrentAccountTransactions200JSONResponse{
		Transactions: &transactions,
		Limit:        &limit,
		Page:         &page,
		MaxPage:      &maxPage,
	}.VisitGetCurrentAccountTransactionsResponse(c.Response())
	return nil
}

// (DELETE /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) MarkDeleteTransactionId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

	// Get transaction from database
	err := s.DBackend.MarkDeleteTransaction(c.Request().Context(), transactionId.String(), adminID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorTransactionNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Transaction %s marked as deleted by %s", transactionId.String(), adminID)
	autogen.MarkDeleteTransactionId200JSONResponse{}.VisitMarkDeleteTransactionIdResponse(c.Response())
	return nil
}

// (GET /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) GetTransactionId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	// Get transaction from database
	transaction, err := s.DBackend.GetTransaction(c.Request().Context(), transactionId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorTransactionNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	autogen.GetTransactionId200JSONResponse(transaction.Transaction).VisitGetTransactionIdResponse(c.Response())
	return nil
}

// (PATCH /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) PatchTransactionId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID, params autogen.PatchTransactionIdParams) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
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
	transaction.State = params.State

	if oldState != autogen.TransactionCanceled && params.State == autogen.TransactionCanceled {
		_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
			err = s.DBackend.UpdateTransaction(ctx, transaction)
			if err != nil {
				logrus.Error(err)
				return nil, errors.New("failed to create transaction")
			}

			// update account balance
			account.Balance += int64(transaction.TotalCost)
			err = s.DBackend.UpdateAccount(ctx, account)
			if err != nil {
				logrus.Error(err)
				return nil, errors.New("failed to update account")
			}

			// update items
			for _, txitem := range transaction.Items {
				item, err := s.DBackend.GetItem(ctx, txitem.ItemId.String())
				if err != nil {
					continue
				}

				item.AmountLeft += txitem.ItemAmount

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
			return Error500(c)
		}
	} else if oldState == autogen.TransactionCanceled && params.State != autogen.TransactionCanceled {
		_, err = s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
			err = s.DBackend.UpdateTransaction(ctx, transaction)
			if err != nil {
				logrus.Error(err)
				return nil, errors.New("failed to create transaction")
			}

			// update account balance
			account.Balance -= int64(transaction.TotalCost)
			err = s.DBackend.UpdateAccount(ctx, account)
			if err != nil {
				logrus.Error(err)
				return nil, errors.New("failed to update account")
			}

			// update items
			for _, txitem := range transaction.Items {
				item, err := s.DBackend.GetItem(ctx, txitem.ItemId.String())
				if err != nil {
					continue
				}

				if item.AmountLeft < txitem.ItemAmount {
					return nil, errors.New("not enough items")
				}

				item.AmountLeft -= txitem.ItemAmount

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
			return Error500(c)
		}
	} else {
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
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
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
		item.ItemAmount = *params.Amount
		item.TotalCost = *params.Amount * item.UnitCost
	}

	origItem, err := s.DBackend.GetItem(c.Request().Context(), itemId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorItemNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	if oldState != autogen.TransactionItemCanceled && item.State == autogen.TransactionItemCanceled {
		origItem.AmountLeft += item.ItemAmount
		account.Balance += int64(item.TotalCost)
		transaction.TotalCost -= item.TotalCost
	} else if oldState == autogen.TransactionItemCanceled && item.State != autogen.TransactionItemCanceled {
		origItem.AmountLeft -= item.ItemAmount
		account.Balance -= int64(item.TotalCost)
		transaction.TotalCost += item.TotalCost
	} else {
		origItem.AmountLeft += oldAmount - item.ItemAmount
		account.Balance += int64(oldCost - item.TotalCost)
		transaction.TotalCost += item.TotalCost - oldCost
	}

	s.DBackend.WithTransaction(c.Request().Context(), func(ctx mongo.SessionContext) (interface{}, error) {
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

	return nil
}

// (GET /transactions)
func (s *Server) GetTransactions(c echo.Context, params autogen.GetTransactionsParams) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	var page uint64
	if params.Page != nil {
		page = *params.Page
	}
	if page > 0 {
		page--
	}

	var limit uint64 = 10
	if params.Limit != nil {
		limit = *params.Limit
	}
	if limit > 100 {
		limit = 100
	}

	var state string
	if params.State != nil {
		state = string(*params.State)
	}

	count, err := s.DBackend.CountAllTransactions(c.Request().Context(), state)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	maxPage := count / limit

	if page > maxPage {
		page = maxPage
	}

	data, err := s.DBackend.GetAllTransactions(c.Request().Context(), page, limit, state)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	transactions := make([]autogen.Transaction, len(data))
	for i, transaction := range data {
		transactions[i] = transaction.Transaction
	}

	autogen.GetTransactions200JSONResponse{
		Transactions: &transactions,
		Limit:        &limit,
		Page:         &page,
		MaxPage:      &maxPage,
	}.VisitGetTransactionsResponse(c.Response())
	return nil
}

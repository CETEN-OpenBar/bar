package api

import (
	"bar/autogen"
	"bar/internal/models"
	"errors"

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
			AccountId: accountID,
			State:     autogen.TransactionStarted,
			Id:        uuid.New(),
		},
	}

	var potentialTransaction autogen.NewTransaction
	var fetchedItems = make(map[string]*models.Item)

	// Check that pin matches
	err = c.Bind(&potentialTransaction)
	if err != nil {
		logrus.Error(err)
		return Error400(c)
	}

	if !account.VerifyPin(potentialTransaction.CardPin) {
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
		} else {
			logrus.Warnf("Item %s already fetched", item.Id.String())
			return Error400(c)
		}

		if item.State == autogen.ItemNotBuyable {
			logrus.Warnf("Item %s is not buyable", item.Id.String())
			return Error400(c)
		}
		if item.AmountLeft < potentialItem.Amount {
			logrus.Warnf("Item %s is not in stock", item.Id.String())
			return Error400(c)
		}
		if item.BuyLimit != nil {
			if *item.BuyLimit < potentialItem.Amount {
				logrus.Warnf("Item %s cannot be bought for that amount", item.Id.String())
				return Error400(c)
			}
		}

		transaction.Items = append(transaction.Items, autogen.TransactionItem{
			ItemAmount: potentialItem.Amount,
			ItemId:     potentialItem.ItemId,
			ItemName:   item.Name,
			PictureUri: item.PictureUri,
			State:      autogen.TransactionItemStarted,
			UnitCost:   item.RealPrice(account.PriceRole),
			TotalCost:  item.RealPrice(account.PriceRole) * potentialItem.Amount,
		})

		transactionCost += item.RealPrice(account.PriceRole) * potentialItem.Amount
		item.AmountLeft -= potentialItem.Amount
	}

	transaction.TotalCost = transactionCost

	if account.Balance < int64(transactionCost) {
		logrus.Warnf("Account %s does not have enough money", accountID)
		return Error400(c)
	}

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
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
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

	maxPage++
	page++
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
	_, err := MustGetUser(c)
	if err != nil {
		return nil
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

	maxPage++
	page++
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
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	// Get transaction from database
	err = s.DBackend.MarkDeleteTransaction(c.Request().Context(), transactionId.String(), account.Id.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorTransactionNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Transaction %s marked as deleted by %s", transactionId.String(), account.Id)
	autogen.MarkDeleteTransactionId200JSONResponse{}.VisitMarkDeleteTransactionIdResponse(c.Response())
	return nil
}

// (GET /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) GetTransactionId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID) error {
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
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

	if oldState != autogen.TransactionItemCanceled && item.State == autogen.TransactionItemCanceled {
		origItem.AmountLeft += item.ItemAmount
		account.Balance += int64(item.TotalCost)
		transaction.TotalCost -= item.TotalCost
	} else if oldState == autogen.TransactionItemCanceled && item.State != autogen.TransactionItemCanceled {
		// Can't do that
		return Error400(c)
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
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
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

	page++
	autogen.GetTransactions200JSONResponse{
		Transactions: &transactions,
		Limit:        &limit,
		Page:         &page,
		MaxPage:      &maxPage,
	}.VisitGetTransactionsResponse(c.Response())
	return nil
}

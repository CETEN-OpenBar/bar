package api

import (
	"bar/autogen"
	"bar/internal/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (POST /account/transactions)
func (s *Server) PostTransactions(c echo.Context) error {
	// Get account from cookie
	sess := s.getUserSess(c)
	accountID, ok := sess.Values["account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(accountID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	transaction := &models.Transaction{
		Transaction: autogen.Transaction{
			AccountId: accountID,
			State:     autogen.TransactionStarted,
			Id:        uuid.New(),
		},
	}

	var potentialTransaction autogen.NewTransaction
	var transactionCost uint64
	for _, potentialItem := range potentialTransaction.Items {
		// Verify that item exists, can be bought, is in stock, and can be bought for that amount
		item, err := s.DBackend.GetItem(potentialItem.ItemId.String())
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return ErrorItemNotFound(c)
			}
			logrus.Error(err)
			return Error500(c)
		}

		if item.State == autogen.ItemNotBuyable {
			return Error400(c)
		}
		if item.AmountLeft < potentialItem.Amount {
			return Error400(c)
		}
		if item.BuyLimit < potentialItem.Amount {
			return Error400(c)
		}

		transaction.Items = append(transaction.Items, autogen.TransactionItem{
			ItemAmount: potentialItem.Amount,
			ItemId:     potentialItem.ItemId,
			State:      autogen.TransactionItemStarted,
			UnitCost:   item.Price,
			TotalCost:  item.Price * potentialItem.Amount,
		})

		transactionCost += item.Price * potentialItem.Amount
	}

	transaction.TotalCost = transactionCost

	err = s.DBackend.CreateTransaction(transaction)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	return nil
}

// (GET /accounts/{account_id}/transactions)
func (s *Server) GetAccountTransactions(c echo.Context, accountId autogen.UUID, params autogen.GetAccountTransactionsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(adminId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
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

	count, err := s.DBackend.CountTransactions(accountId.String(), state)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	maxPage := count / limit

	if page > maxPage {
		page = maxPage
	}

	data, err := s.DBackend.GetTransactions(accountId.String(), page, limit, state)
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
	sess := s.getUserSess(c)
	accountID, ok := sess.Values["account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(accountID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
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

	count, err := s.DBackend.CountTransactions(accountID, state)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	maxPage := count / limit

	if page > maxPage {
		page = maxPage
	}

	data, err := s.DBackend.GetTransactions(accountID, page, limit, state)
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
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(adminId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	// Get transaction from database
	err = s.DBackend.MarkDeleteTransaction(transactionId.String(), adminId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorTransactionNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Transaction %s marked as deleted by %s", transactionId.String(), adminId)
	autogen.MarkDeleteTransactionId200JSONResponse{}.VisitMarkDeleteTransactionIdResponse(c.Response())
	return nil
}

// (GET /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) GetTransactionId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(adminId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	// Get transaction from database
	transaction, err := s.DBackend.GetTransaction(transactionId.String())
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
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(adminId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	// Get transaction from database
	transaction, err := s.DBackend.GetTransaction(transactionId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorTransactionNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	transaction.State = params.State

	err = s.DBackend.UpdateTransaction(transaction)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	return nil
}

// (PATCH /accounts/{account_id}/transactions/{transaction_id}/{item_id})
func (s *Server) PatchTransactionItemId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID, itemId autogen.UUID, params autogen.PatchTransactionItemIdParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(adminId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	// Get transaction from database
	transaction, err := s.DBackend.GetTransaction(transactionId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorTransactionNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	for i, item := range transaction.Items {
		if item.ItemId == itemId {
			if params.State != nil {
				transaction.Items[i].State = *params.State
			}
			if params.Amount != nil {
				transaction.Items[i].ItemAmount = *params.Amount
				transaction.Items[i].TotalCost = *params.Amount * transaction.Items[i].UnitCost
			}

			break
		}
	}

	err = s.DBackend.UpdateTransaction(transaction)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}
	return nil
}

// (GET /transactions)
func (s *Server) GetTransactions(c echo.Context, params autogen.GetTransactionsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(adminId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		return Error500(c)
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

	count, err := s.DBackend.CountAllTransactions(state)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	maxPage := count / limit

	if page > maxPage {
		page = maxPage
	}

	data, err := s.DBackend.GetAllTransactions(page, limit, state)
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

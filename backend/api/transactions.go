package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

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

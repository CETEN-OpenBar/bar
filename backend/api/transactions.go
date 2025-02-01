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

	var state string
	if params.State != nil {
		state = string(*params.State)
	}

	count, err := s.DBackend.CountTransactions(c.Request().Context(), accountId.String(), state)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetTransactions(c.Request().Context(), accountId.String(), dbpage, limit, state)
	if err != nil {
		return Error500(c)
	}

	transactions := make([]autogen.Transaction, len(data))
	for i, transaction := range data {
		transactions[i] = transaction.Transaction
	}

	autogen.GetAccountTransactions200JSONResponse{
		Transactions: transactions,
		Limit:        limit,
		Page:         page,
		MaxPage:      maxPage,
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

	var state string
	if params.State != nil {
		state = string(*params.State)
	}

	count, err := s.DBackend.CountTransactions(c.Request().Context(), accountID, state)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetTransactions(c.Request().Context(), accountID, dbpage, limit, state)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	transactions := make([]autogen.Transaction, len(data))
	for i, transaction := range data {
		transactions[i] = transaction.Transaction
	}

	autogen.GetCurrentAccountTransactions200JSONResponse{
		Transactions: transactions,
		Limit:        limit,
		Page:         page,
		MaxPage:      maxPage,
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

	logrus.WithField("transaction", transactionId.String()).WithField("by", account.Name()).Info("Transaction marked for deletion")
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

	var state string
	if params.State != nil {
		state = string(*params.State)
	}

	var name string
	if params.Name != nil {
		name = string(*params.Name)
	}

	var hide_remotes bool
	if params.HideRemote != nil {
		hide_remotes = bool(*params.HideRemote)
	} else {
		hide_remotes = true
	}

	count, err := s.DBackend.CountAllTransactions(c.Request().Context(), state, name, hide_remotes)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetAllTransactions(c.Request().Context(), dbpage, limit, state, name, hide_remotes)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	transactions := make([]autogen.Transaction, len(data))
	for i, transaction := range data {
		transactions[i] = transaction.Transaction
	}

	autogen.GetTransactions200JSONResponse{
		Transactions: transactions,
		Limit:        limit,
		Page:         page,
		MaxPage:      maxPage,
	}.VisitGetTransactionsResponse(c.Response())
	return nil
}


// (GET /transactions/items)
func (s *Server) GetTransactionsItems(c echo.Context, params autogen.GetTransactionsItemsParams) error {
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	var name string
	if params.Name != nil {
		name = string(*params.Name)
	}

	data, err := s.DBackend.GetAllActiveTransactionsItems(c.Request().Context(), name)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	autogen.GetTransactionsItems200JSONResponse(data).VisitGetTransactionsItemsResponse(c.Response())
	return nil
}


func (s *Server) GetTransactionsByTimestamp(c echo.Context, params autogen.GetTransactionsByTimestampParams) error {
	_, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}
	
	StartTime := uint64(params.StartTime)
	EndTime := uint64(params.EndTime)

	if StartTime > EndTime {
		return Error400(c)
	}

	count, err := s.DBackend.CountTransactionsByTimestamp(c.Request().Context(), StartTime, EndTime)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetTransactionsByTimestamp(c.Request().Context(), StartTime, EndTime, dbpage, limit)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	transactions := make([]autogen.Transaction, len(data))
	for i, transaction := range data {
		transactions[i] = transaction.Transaction
	}

	autogen.GetTransactionsByTimestamp200JSONResponse{
		Transactions: transactions,
		Limit:        limit,
		Page:         page,
		MaxPage:      maxPage,
	}.VisitGetTransactionsByTimestampResponse(c.Response())
	return nil
}

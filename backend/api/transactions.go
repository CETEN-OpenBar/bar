package api

import (
	"bar/autogen"
	"bar/internal/models"

	"github.com/labstack/echo/v4"
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
		return Error500(c)
	}

	var transaction models.Transaction
	var potentialTransaction autogen.NewTransaction
	var transactionCost uint64
	for _, potentialItem := range potentialTransaction.Items {
		// Verify that item exists, can be bought, is in stock, and can be bought for that amount
		item, err := s.DBackend.GetItem(potentialItem.ItemId.String())
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return ErrorItemNotFound(c)
			}
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
			State:      autogen.TransactionStarted,
			TotalCost:  item.Price * potentialItem.Amount,
		})

		transactionCost += item.Price * potentialItem.Amount

	}

	return nil
}

// (GET /accounts/{account_id}/transactions)
func (s *Server) GetAccountTransactions(c echo.Context, accountId autogen.UUID, params autogen.GetAccountTransactionsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
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
		return Error500(c)
	}

	// TODO: implement
	return nil
}

// (DELETE /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) MarkDeleteTransactionId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (GET /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) GetTransactionId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /accounts/{account_id}/transactions/{transaction_id})
func (s *Server) PatchTransactionId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID, params autogen.PatchTransactionIdParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /accounts/{account_id}/transactions/{transaction_id}/{item_id})
func (s *Server) PatchTransactionItemId(c echo.Context, accountId autogen.UUID, transactionId autogen.UUID, itemId autogen.UUID, params autogen.PatchTransactionItemIdParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (GET /transactions)
func (s *Server) GetTransactions(c echo.Context, params autogen.GetTransactionsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

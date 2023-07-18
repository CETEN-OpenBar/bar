package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
)

// (GET /deleted/accounts)
func (s *Server) GetDeletedAccounts(c echo.Context, params autogen.GetDeletedAccountsParams) error {
	// TODO: implement
	return nil
}

// (DELETE /deleted/accounts/{account_id})
func (s *Server) DeleteAccount(c echo.Context, accountId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /deleted/accounts/{account_id})
func (s *Server) RestoreDeletedAccount(c echo.Context, accountId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (GET /deleted/carousel/images)
func (s *Server) GetDeletedCarouselImages(c echo.Context, params autogen.GetDeletedCarouselImagesParams) error {
	// TODO: implement
	return nil
}

// (DELETE /deleted/carousel/images/{image_id})
func (s *Server) DeleteCarouselImage(c echo.Context, imageId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /deleted/carousel/images/{image_id})
func (s *Server) RestoreDeletedCarouselImage(c echo.Context, imageId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (GET /deleted/carousel/texts)
func (s *Server) GetDeletedCarouselTexts(c echo.Context, params autogen.GetDeletedCarouselTextsParams) error {
	// TODO: implement
	return nil
}

// (DELETE /deleted/carousel/texts/{text_id})
func (s *Server) DeleteCarouselText(c echo.Context, textId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /deleted/carousel/texts/{text_id})
func (s *Server) RestoreDeletedCarouselText(c echo.Context, textId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (GET /deleted/items)
func (s *Server) GetDeletedItems(c echo.Context, params autogen.GetDeletedItemsParams) error {
	// TODO: implement
	return nil
}

// (DELETE /deleted/items/{item_id})
func (s *Server) DeleteItem(c echo.Context, itemId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /deleted/items/{item_id})
func (s *Server) RestoreDeletedItem(c echo.Context, itemId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (GET /deleted/refills)
func (s *Server) GetDeletedRefills(c echo.Context, params autogen.GetDeletedRefillsParams) error {
	// TODO: implement
	return nil
}

// (DELETE /deleted/refills/{refill_id})
func (s *Server) DeleteRefill(c echo.Context, refillId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /deleted/refills/{refill_id})
func (s *Server) RestoreDeletedRefill(c echo.Context, refillId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (GET /deleted/transactions)
func (s *Server) GetDeletedTransactions(c echo.Context, params autogen.GetDeletedTransactionsParams) error {
	// TODO: implement
	return nil
}

// (DELETE /deleted/transactions/{transaction_id})
func (s *Server) DeleteTransaction(c echo.Context, transactionId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /deleted/transactions/{transaction_id})
func (s *Server) RestoreDeletedTransaction(c echo.Context, transactionId autogen.UUID) error {
	// TODO: implement
	return nil
}

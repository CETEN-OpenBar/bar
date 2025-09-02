package api

import (
	"bar/autogen"
	"bar/internal/storage"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (GET /deleted/accounts)
func (s *Server) GetDeletedAccounts(c echo.Context, params autogen.GetDeletedAccountsParams) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	var search string
	if params.Search != nil {
		search = *params.Search
	}

	count, err := s.DBackend.CountDeletedAccounts(c.Request().Context(), search)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetDeletedAccounts(c.Request().Context(), dbpage, limit, search)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var accounts []autogen.Account

	for _, acc := range data {
		accounts = append(accounts, acc.Account)
	}

	autogen.GetDeletedAccounts200JSONResponse{
		Accounts: accounts,
		Limit:    limit,
		Page:     page,
		MaxPage:  maxPage,
	}.VisitGetDeletedAccountsResponse(c.Response())
	return nil
}

// (DELETE /deleted/accounts/{account_id})
func (s *Server) DeleteAccount(c echo.Context, accountId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.DeleteAccount(c.Request().Context(), accountId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("account", accountId.String()).WithField("by", account.Name()).Info("Account deleted")
	return nil
}

// (PATCH /deleted/accounts/{account_id})
func (s *Server) RestoreDeletedAccount(c echo.Context, accountId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.UnMarkDeleteAccount(c.Request().Context(), accountId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorAccNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}
	logrus.WithField("account", accountId.String()).WithField("by", account.Name()).Info("Account restored")
	return nil
}

// (GET /deleted/carousel/images)
func (s *Server) GetDeletedCarouselImages(c echo.Context, params autogen.GetDeletedCarouselImagesParams) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	count, err := s.DBackend.CountDeletedCarouselImages(c.Request().Context())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetDeletedCarouselImages(c.Request().Context(), dbpage, limit)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var images []autogen.CarouselImage

	for _, acc := range data {
		images = append(images, acc.CarouselImage)
	}

	autogen.GetDeletedCarouselImages200JSONResponse{
		Items:   images,
		Limit:   limit,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetDeletedCarouselImagesResponse(c.Response())
	return nil
}

// (DELETE /deleted/carousel/images/{image_id})
func (s *Server) DeleteCarouselImage(c echo.Context, imageId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.DeleteCarouselImage(c.Request().Context(), imageId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorImageNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	err = storage.DeleteFile("carousel/" + imageId.String())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("image", imageId.String()).WithField("by", account.Name()).Info("Carousel image deleted")
	return nil
}

// (PATCH /deleted/carousel/images/{image_id})
func (s *Server) RestoreDeletedCarouselImage(c echo.Context, imageId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.UnMarkDeleteCarouselImage(c.Request().Context(), imageId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorImageNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}
	logrus.WithField("image", imageId.String()).WithField("by", account.Name()).Info("Carousel image restored")
	return nil
}

// (GET /deleted/carousel/texts)
func (s *Server) GetDeletedCarouselTexts(c echo.Context, params autogen.GetDeletedCarouselTextsParams) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	count, err := s.DBackend.CountDeletedCarouselTexts(c.Request().Context())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetDeletedCarouselTexts(c.Request().Context(), dbpage, limit)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var items []autogen.CarouselText

	for _, acc := range data {
		items = append(items, acc.CarouselText)
	}

	autogen.GetDeletedCarouselTexts200JSONResponse{
		Items:   items,
		Limit:   limit,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetDeletedCarouselTextsResponse(c.Response())
	return nil
}

// (DELETE /deleted/carousel/texts/{text_id})
func (s *Server) DeleteCarouselText(c echo.Context, textId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.DeleteCarouselText(c.Request().Context(), textId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorTextNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("text", textId.String()).WithField("by", account.Name()).Info("Carousel text deleted")
	return nil
}

// (PATCH /deleted/carousel/texts/{text_id})
func (s *Server) RestoreDeletedCarouselText(c echo.Context, textId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.UnMarkDeleteCarouselText(c.Request().Context(), textId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorTextNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}
	return nil
}

// (GET /deleted/categories)
func (s *Server) GetDeletedCategories(c echo.Context, params autogen.GetDeletedCategoriesParams) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	count, err := s.DBackend.CountDeletedCategories(c.Request().Context())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetDeletedCategories(c.Request().Context(), dbpage, limit)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var items []autogen.Category

	for _, acc := range data {
		items = append(items, acc.Category)
	}

	autogen.GetDeletedCategories200JSONResponse{
		Categories: items,
		Limit:      limit,
		Page:       page,
		MaxPage:    maxPage,
	}.VisitGetDeletedCategoriesResponse(c.Response())
	return nil
}

// (DELETE /deleted/categories/{category_id})
func (s *Server) DeleteCategory(c echo.Context, categoryId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.DeleteCategory(c.Request().Context(), categoryId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorItemNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	err = storage.DeleteFile("categories/" + categoryId.String())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("category", categoryId.String()).WithField("by", account.Name()).Info("Category deleted")
	return nil
}

// (PATCH /deleted/categories/{category_id})
func (s *Server) RestoreDeletedCategory(c echo.Context, categoryId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.UnMarkDeleteCategory(c.Request().Context(), categoryId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorItemNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}
	logrus.WithField("category", categoryId.String()).WithField("by", account.Name()).Info("Category restored")
	return nil
}

// (GET /deleted/items)
func (s *Server) GetDeletedItems(c echo.Context, params autogen.GetDeletedItemsParams) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	count, err := s.DBackend.CountDeletedItems(c.Request().Context())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetDeletedItems(c.Request().Context(), dbpage, limit)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var items []autogen.Item

	for _, acc := range data {
		items = append(items, acc.Item)
	}

	autogen.GetDeletedItems200JSONResponse{
		Items:   items,
		Limit:   limit,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetDeletedItemsResponse(c.Response())
	return nil
}

// (DELETE /deleted/items/{item_id})
func (s *Server) DeleteItem(c echo.Context, itemId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.DeleteItem(c.Request().Context(), itemId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorItemNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	err = storage.DeleteFile("items/" + itemId.String())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("item", itemId.String()).WithField("by", account.Name()).Info("Item deleted")
	return nil
}

// (PATCH /deleted/items/{item_id})
func (s *Server) RestoreDeletedItem(c echo.Context, itemId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.UnMarkDeleteItem(c.Request().Context(), itemId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorItemNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}
	logrus.WithField("item", itemId.String()).WithField("by", account.Name()).Info("Item restored")
	return nil
}

// (GET /deleted/refills)
func (s *Server) GetDeletedRefills(c echo.Context, params autogen.GetDeletedRefillsParams) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	count, err := s.DBackend.CountDeletedRefills(c.Request().Context())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetDeletedRefills(c.Request().Context(), dbpage, limit)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var items []autogen.Refill

	for _, acc := range data {
		items = append(items, acc.Refill)
	}

	autogen.GetDeletedRefills200JSONResponse{
		Refills: items,
		Limit:   limit,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetDeletedRefillsResponse(c.Response())
	return nil
}

// (DELETE /deleted/refills/{refill_id})
func (s *Server) DeleteRefill(c echo.Context, refillId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.DeleteRefill(c.Request().Context(), refillId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorRefillNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("refill", refillId.String()).WithField("by", account.Name()).Info("Refill deleted")
	return nil
}

// (GET /deleted/refills)
func (s *Server) GetDeletedStarring(c echo.Context, params autogen.GetDeletedStarringParams) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	count, err := s.DBackend.CountDeletedStarrings(c.Request().Context())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetDeletedStarrings(c.Request().Context(), dbpage, limit)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var items []autogen.Starring

	for _, acc := range data {
		items = append(items, acc.Starring)
	}

	autogen.GetDeletedStarring200JSONResponse{
		Starring: items,
		Limit:   limit,
		Page:    page,
		MaxPage: maxPage,
	}.VisitGetDeletedStarringResponse(c.Response())
	return nil
}

// (DELETE /deleted/stars/{starring_id})
func (s *Server) DeleteStarring(c echo.Context, starringId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.DeleteStarring(c.Request().Context(), starringId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorRefillNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("starring", starringId.String()).WithField("by", account.Name()).Info("Starring deleted")
	return nil
}

// (PATCH /deleted/refills/{refill_id})
func (s *Server) RestoreDeletedRefill(c echo.Context, refillId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.UnMarkDeleteRefill(c.Request().Context(), refillId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorRefillNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}
	logrus.WithField("refill", refillId.String()).WithField("by", account.Name()).Info("Refill restored")
	return nil
}

// (PATCH /deleted/refills/{refill_id})
func (s *Server) RestoreDeletedStarring(c echo.Context, starringId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.UnMarkDeleteStarring(c.Request().Context(), starringId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorStarringNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}
	logrus.WithField("starring", starringId.String()).WithField("by", account.Name()).Info("Starring restored")
	return nil
}

// (GET /deleted/transactions)
func (s *Server) GetDeletedTransactions(c echo.Context, params autogen.GetDeletedTransactionsParams) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	count, err := s.DBackend.CountDeletedTransactions(c.Request().Context())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetDeletedTransactions(c.Request().Context(), dbpage, limit)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var items []autogen.Transaction

	for _, acc := range data {
		items = append(items, acc.Transaction)
	}

	autogen.GetDeletedTransactions200JSONResponse{
		Transactions: items,
		Limit:        limit,
		Page:         page,
		MaxPage:      maxPage,
	}.VisitGetDeletedTransactionsResponse(c.Response())
	return nil
}

// (DELETE /deleted/transactions/{transaction_id})
func (s *Server) DeleteTransaction(c echo.Context, transactionId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.DeleteTransaction(c.Request().Context(), transactionId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorTransactionNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("transaction", transactionId.String()).WithField("by", account.Name()).Info("Transaction deleted")
	return nil
}

// (PATCH /deleted/transactions/{transaction_id})
func (s *Server) RestoreDeletedTransaction(c echo.Context, transactionId autogen.UUID) error {
	// Get admin account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	if account.Role != autogen.AccountSuperAdmin {
		return ErrorNotAuthenticated(c)
	}

	err = s.DBackend.UnMarkDeleteTransaction(c.Request().Context(), transactionId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorTransactionNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}
	logrus.WithField("transaction", transactionId.String()).WithField("by", account.Name()).Info("Transaction restored")
	return nil
}

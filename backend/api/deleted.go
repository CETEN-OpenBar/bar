package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// (GET /deleted/accounts)
func (s *Server) GetDeletedAccounts(c echo.Context, params autogen.GetDeletedAccountsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	var page int
	if params.Page != nil {
		page = *params.Page
	}

	var limit = 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	count, err := s.DBackend.CountDeletedAccounts()
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var maxPage = int(count) / limit

	if page > maxPage {
		page = maxPage
	}

	if page < 0 {
		page = 0
	}

	if limit < 0 {
		limit = 0
	}

	if limit > 100 {
		limit = 100
	}

	data, err := s.DBackend.GetDeletedAccounts(page, limit)
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
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	err := s.DBackend.DeleteAccount(accountId.String())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	adminId, _ := sess.Values["admin_account_id"].(string)
	logrus.Infof("Admin %s deleted account %s", adminId, accountId)
	return nil
}

// (PATCH /deleted/accounts/{account_id})
func (s *Server) RestoreDeletedAccount(c echo.Context, accountId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	err := s.DBackend.UnMarkDeleteAccount(accountId.String())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}
	return nil
}

// (GET /deleted/carousel/images)
func (s *Server) GetDeletedCarouselImages(c echo.Context, params autogen.GetDeletedCarouselImagesParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	var page int
	if params.Page != nil {
		page = *params.Page
	}

	var limit = 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	count, err := s.DBackend.CountDeletedCarouselImages()
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var maxPage = int(count) / limit

	if page > maxPage {
		page = maxPage
	}

	if page < 0 {
		page = 0
	}

	if limit < 0 {
		limit = 0
	}

	if limit > 100 {
		limit = 100
	}

	data, err := s.DBackend.GetDeletedCarouselImages(page, limit)
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
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	err := s.DBackend.DeleteCarouselImage(imageId.String())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	adminId, _ := sess.Values["admin_account_id"].(string)
	logrus.Infof("Admin %s deleted carousel image %s", adminId, imageId)
	return nil
}

// (PATCH /deleted/carousel/images/{image_id})
func (s *Server) RestoreDeletedCarouselImage(c echo.Context, imageId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	return nil
}

// (GET /deleted/carousel/texts)
func (s *Server) GetDeletedCarouselTexts(c echo.Context, params autogen.GetDeletedCarouselTextsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	var page int
	if params.Page != nil {
		page = *params.Page
	}

	var limit = 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	count, err := s.DBackend.CountDeletedCarouselTexts()
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var maxPage = int(count) / limit

	if page > maxPage {
		page = maxPage
	}

	if page < 0 {
		page = 0
	}

	if limit < 0 {
		limit = 0
	}

	if limit > 100 {
		limit = 100
	}

	data, err := s.DBackend.GetDeletedCarouselTexts(page, limit)
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
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /deleted/carousel/texts/{text_id})
func (s *Server) RestoreDeletedCarouselText(c echo.Context, textId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (GET /deleted/items)
func (s *Server) GetDeletedItems(c echo.Context, params autogen.GetDeletedItemsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	var page int
	if params.Page != nil {
		page = *params.Page
	}

	var limit = 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	count, err := s.DBackend.CountDeletedItems()
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var maxPage = int(count) / limit

	if page > maxPage {
		page = maxPage
	}

	if page < 0 {
		page = 0
	}

	if limit < 0 {
		limit = 0
	}

	if limit > 100 {
		limit = 100
	}

	data, err := s.DBackend.GetDeletedItems(page, limit)
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
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /deleted/items/{item_id})
func (s *Server) RestoreDeletedItem(c echo.Context, itemId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (GET /deleted/refills)
func (s *Server) GetDeletedRefills(c echo.Context, params autogen.GetDeletedRefillsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	var page int
	if params.Page != nil {
		page = *params.Page
	}

	var limit = 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	count, err := s.DBackend.CountDeletedRefills()
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var maxPage = int(count) / limit

	if page > maxPage {
		page = maxPage
	}

	if page < 0 {
		page = 0
	}

	if limit < 0 {
		limit = 0
	}

	if limit > 100 {
		limit = 100
	}

	data, err := s.DBackend.GetDeletedRefills(page, limit)
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
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /deleted/refills/{refill_id})
func (s *Server) RestoreDeletedRefill(c echo.Context, refillId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (GET /deleted/transactions)
func (s *Server) GetDeletedTransactions(c echo.Context, params autogen.GetDeletedTransactionsParams) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	var page int
	if params.Page != nil {
		page = *params.Page
	}

	var limit = 10
	if params.Limit != nil {
		limit = *params.Limit
	}

	count, err := s.DBackend.CountDeletedTransactions()
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var maxPage = int(count) / limit

	if page > maxPage {
		page = maxPage
	}

	if page < 0 {
		page = 0
	}

	if limit < 0 {
		limit = 0
	}

	if limit > 100 {
		limit = 100
	}

	data, err := s.DBackend.GetDeletedTransactions(page, limit)
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
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /deleted/transactions/{transaction_id})
func (s *Server) RestoreDeletedTransaction(c echo.Context, transactionId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["super_admin"].(bool)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

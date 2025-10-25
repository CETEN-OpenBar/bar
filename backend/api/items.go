package api

import (
	"bar/autogen"
	"bar/internal/storage"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (GET /categories/{category_id}/items)
func (s *Server) GetCategoryItems(c echo.Context, categoryId autogen.UUID, params autogen.GetCategoryItemsParams) error {
	// Get account from cookie
	account, err := MustGetUser(c)
	if err != nil {
		return nil
	}
	state := ""
	if params.State != nil {
		state = string(*params.State)
	}

	_, err = s.DBackend.GetCategory(c.Request().Context(), categoryId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorCategoryNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	count, err := s.DBackend.CountItems(c.Request().Context(), categoryId.String(), state, "", "", "")
	if err != nil {
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetItems(c.Request().Context(), categoryId.String(), dbpage, limit, state, "", "", "", true)
	if err != nil {
		return Error500(c)
	}

	var items []autogen.Item

	for _, item := range data {
		rp := item.RealPrice(account.PriceRole)
		item.DisplayPrice = &rp

		if account.HasPrivileges() {
			rp := item.RealPrices()
			item.DisplayPrices = &rp
		}

		items = append(items, item.Item)
	}

	autogen.GetCategoryItems200JSONResponse{
		Items:   items,
		Page:    page,
		Limit:   limit,
		MaxPage: maxPage,
	}.VisitGetCategoryItemsResponse(c.Response())

	return nil
}

// (DELETE /categories/{category_id}/items/{item_id})
func (s *Server) MarkDeleteItem(c echo.Context, categoryId autogen.UUID, itemId autogen.UUID) error {
	admin, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	_, err = s.DBackend.GetItem(c.Request().Context(), itemId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorItemNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	err = s.DBackend.MarkDeleteItem(c.Request().Context(), itemId.String(), admin.Id.String())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("item", itemId.String()).WithField("by", admin.Name()).Info("Item marked for deletion")
	autogen.DeleteItem204Response{}.VisitDeleteItemResponse(c.Response())
	return nil
}

// (PATCH /categories/{category_id}/items/{item_id})
func (s *Server) PatchItem(c echo.Context, categoryId autogen.UUID, itemId autogen.UUID) error {
	admin, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	item, err := s.DBackend.GetItem(c.Request().Context(), itemId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorItemNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	var p autogen.UpdateItem
	if err := c.Bind(&p); err != nil {
		return Error400(c)
	}

	if p.Picture != nil {
		d, err := base64.StdEncoding.DecodeString(*p.Picture)
		if err != nil {
			return Error400(c)
		}

		// Check MIME type
		if !strings.Contains(http.DetectContentType(d), "image") {
			return Error400(c)
		}

		// Save image to storage
		err = storage.SaveFile("items/"+item.Id.String(), d)
		if err != nil {
			logrus.Error(err)
			return Error500(c)
		}
	}

	if p.CategoryId != nil {
		item.CategoryId = *p.CategoryId
	}
	if p.Name != nil {
		item.Name = *p.Name
	}
	if p.Prices != nil {
		item.Prices = *p.Prices
	}

	if p.Promotion != nil {
		item.Promotion = p.Promotion
	}

	if p.PromotionEndsAt != nil {
		item.PromotionEndsAt = p.PromotionEndsAt
	}

	if p.State != nil {
		item.State = *p.State
	}
	if p.AmountLeft != nil {
		item.AmountLeft = *p.AmountLeft
	}
	if p.OptimalAmount != nil {
		item.OptimalAmount = *p.OptimalAmount
	}
	if p.BuyLimit != nil {
		if *p.BuyLimit < 0 {
			item.BuyLimit = nil
		} else {
			var buyLimit uint64 = uint64(*p.BuyLimit)
			item.BuyLimit = &buyLimit
		}
	}
	if p.AmountPerBundle != nil {
		item.AmountPerBundle = p.AmountPerBundle
	}
	if p.RefBundle != nil {
		item.RefBundle = p.RefBundle
	}
	if p.Fournisseur != nil {
		item.Fournisseur = p.Fournisseur
	}

	rp := item.RealPrices()
	item.DisplayPrices = &rp

	// Save item to database
	err = s.DBackend.UpdateItem(c.Request().Context(), item)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("item", item.Name).WithField("by", admin.Name()).Info("Item updated")
	autogen.PostItem201JSONResponse(item.Item).VisitPostItemResponse(c.Response())
	return nil
}

// (GET /categories/{category_id}/items/{item_id}/picture)
func (s *Server) GetItemPicture(c echo.Context, categoryId autogen.UUID, itemId autogen.UUID) error {
	// Get account from cookie
	_, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	data, err := storage.GetFile("items/" + itemId.String())
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			// Remove cache
			c.Response().Header().Set("Cache-Control", "max-age=0")
			c.Response().Header().Set("Expires", "0")
			return ErrorCategoryNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	c.Response().Header().Set("Cache-Control", "max-age: 0, must-revalidate")
	c.Response().Header().Set("Expires", "0")

	// ETag is sha1 of data
	hash := sha1.Sum(data)
	c.Response().Header().Set("ETag", fmt.Sprintf("%x", hash))
	// Check "If-None-Match" header
	if c.Request().Header.Get("If-None-Match") == fmt.Sprintf("%x", hash) {
		c.Response().WriteHeader(http.StatusNotModified)
		return nil
	}

	c.Response().Header().Set("Content-Type", http.DetectContentType(data))
	c.Response().Write(data)
	return nil
}

// (GET /items)
func (s *Server) GetAllItems(c echo.Context, params autogen.GetAllItemsParams) error {
	// Get account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	state := ""
	categoryId := ""
	name := ""
	fournisseur := ""
	refBundle := ""
	if params.State != nil {
		state = string(*params.State)
	}
	if params.CategoryId != nil {
		categoryId = params.CategoryId.String()
	}
	if params.Name != nil {
		name = string(*params.Name)
	}
	if params.Fournisseur != nil {
		fournisseur = string(*params.Fournisseur)
	}
	if params.RefBundle != nil {
		refBundle = string(*params.RefBundle)
	}

	count, err := s.DBackend.CountItems(c.Request().Context(), categoryId, state, name, fournisseur, refBundle)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetItems(c.Request().Context(), categoryId, dbpage, limit, state, name, fournisseur, refBundle, false)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var items []autogen.Item

	for _, item := range data {
		rp := item.RealPrice(account.PriceRole)
		item.DisplayPrice = &rp

		if account.HasPrivileges() {
			rp := item.RealPrices()
			item.DisplayPrices = &rp
		}

		items = append(items, item.Item)
	}

	autogen.GetAllItems200JSONResponse{
		Items:   items,
		Page:    page,
		Limit:   limit,
		MaxPage: maxPage,
	}.VisitGetAllItemsResponse(c.Response())

	return nil
}

// (GET /items/incoherent)
func (s *Server) GetAllIncoherentItems(c echo.Context, params autogen.GetAllIncoherentItemsParams) error {
	// Get account from cookie
	account, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	state := ""
	categoryId := ""
	name := ""
	if params.State != nil {
		state = string(*params.State)
	}
	if params.CategoryId != nil {
		categoryId = params.CategoryId.String()
	}
	if params.Name != nil {
		name = string(*params.Name)
	}

	count, err := s.DBackend.CountIncoherentItems(c.Request().Context(), categoryId, state, name)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Make sure the last page is not empty
	dbpage, page, limit, maxPage := autogen.Pager(params.Page, params.Limit, &count)

	data, err := s.DBackend.GetIncoherentItems(c.Request().Context(), dbpage, limit, categoryId, state, name)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	

	var items []autogen.Item

	for _, item := range data {
		rp := item.RealPrice(account.PriceRole)
		item.DisplayPrice = &rp

		if account.HasPrivileges() {
			rp := item.RealPrices()
			item.DisplayPrices = &rp
		}

		items = append(items, item.Item)
	}

	autogen.GetAllIncoherentItems200JSONResponse{
		Items:   items,
		Page:    page,
		Limit:   limit,
		MaxPage: maxPage,
	}.VisitGetAllIncoherentItemsResponse(c.Response())

	return nil
}
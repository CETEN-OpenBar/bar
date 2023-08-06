package api

import (
	"bar/autogen"
	"bar/internal/models"
	"bar/internal/storage"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (GET /categories/{category_id}/items)
func (s *Server) GetCategoryItems(c echo.Context, categoryId autogen.UUID, params autogen.GetCategoryItemsParams) error {
	// Get account from cookie
	logged := c.Get("userLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	var page uint64 = 1
	if params.Page != nil {
		page = uint64(*params.Page)
	}

	var size uint64 = 50
	if params.Limit != nil {
		size = uint64(*params.Limit)
	}

	if page > 0 {
		page -= 1
	}

	if size > 100 {
		size = 100
	}

	var state = ""
	if params.State != nil {
		state = string(*params.State)
	}

	_, err := s.DBackend.GetCategory(c.Request().Context(), categoryId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorCategoryNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	count, err := s.DBackend.CountItems(c.Request().Context(), categoryId.String(), state)
	if err != nil {
		return Error500(c)
	}
	var maxPage = uint64(count) / size

	if page > maxPage {
		page = maxPage
	}

	data, err := s.DBackend.GetItems(c.Request().Context(), categoryId.String(), page, size, state)
	if err != nil {
		return Error500(c)
	}

	var items []autogen.Item

	for _, item := range data {
		var rp = item.RealPrice()
		item.DisplayPrice = &rp
		items = append(items, item.Item)
	}

	page += 1
	maxPage += 1
	autogen.GetCategoryItems200JSONResponse{
		Items:   &items,
		Page:    &page,
		Limit:   &size,
		MaxPage: &maxPage,
	}.VisitGetCategoryItemsResponse(c.Response())

	return nil
}

// (POST /categories/{category_id}/items)
func (s *Server) PostItem(c echo.Context, categoryId autogen.UUID) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	var p autogen.NewItem
	if err := c.Bind(&p); err != nil {
		logrus.Error(err)
		return Error400(c)
	}

	d, err := base64.StdEncoding.DecodeString(p.Picture)
	if err != nil {
		logrus.Error(err)
		return Error400(c)
	}

	// Check MIME type
	if !strings.Contains(http.DetectContentType(d), "image") {
		logrus.Error(err)
		return Error400(c)
	}

	uid := uuid.New()

	// Save image to storage
	err = storage.SaveFile("items/"+uid.String(), d)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	item := &models.Item{
		Item: autogen.Item{
			Id:              uid,
			CategoryId:      categoryId,
			Name:            p.Name,
			Price:           p.Price,
			PictureUri:      "/categories/" + categoryId.String() + "/items/" + uid.String() + "/picture",
			Promotion:       p.Promotion,
			PromotionEndsAt: p.PromotionEndsAt,
			State:           p.State,
			AmountLeft:      p.AmountLeft,
			BuyLimit:        p.BuyLimit,
		},
	}

	// Save item to database
	err = s.DBackend.CreateItem(c.Request().Context(), item)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Item %s created by %s", item.Id.String(), item.CategoryId)
	autogen.PostItem201JSONResponse(item.Item).VisitPostItemResponse(c.Response())
	return nil
}

// (DELETE /categories/{category_id}/items/{item_id})
func (s *Server) MarkDeleteItem(c echo.Context, categoryId autogen.UUID, itemId autogen.UUID) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

	_, err := s.DBackend.GetItem(c.Request().Context(), itemId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorItemNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	err = s.DBackend.MarkDeleteItem(c.Request().Context(), itemId.String(), adminID)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Item %s deleted by %s", itemId.String(), adminID)
	autogen.DeleteItem204Response{}.VisitDeleteItemResponse(c.Response())
	return nil
}

// (PATCH /categories/{category_id}/items/{item_id})
func (s *Server) PatchItem(c echo.Context, categoryId autogen.UUID, itemId autogen.UUID) error {
	logged := c.Get("adminLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
	}

	adminID := c.Get("adminAccountID").(string)

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
	if p.Price != nil {
		item.Price = *p.Price
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
	if p.BuyLimit != nil {
		item.BuyLimit = *p.BuyLimit
	}

	var rp = item.RealPrice()
	item.DisplayPrice = &rp

	// Save item to database
	err = s.DBackend.UpdateItem(c.Request().Context(), item)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Item %s updated by %s", item.Id.String(), adminID)
	autogen.PostItem201JSONResponse(item.Item).VisitPostItemResponse(c.Response())
	return nil
}

// (GET /categories/{category_id}/items/{item_id}/picture)
func (s *Server) GetItemPicture(c echo.Context, categoryId autogen.UUID, itemId autogen.UUID) error {
	// Get account from cookie
	logged := c.Get("userLogged").(bool)
	if !logged {
		return ErrorNotAuthenticated(c)
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

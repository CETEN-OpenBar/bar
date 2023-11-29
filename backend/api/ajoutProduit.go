package api

import (
	"bar/autogen"
	"bar/internal/models"
	"bar/internal/storage"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostItem403(c echo.Context, message string) error {
	r := autogen.PostItem403JSONResponse{
		ErrorCode: autogen.ErrForbidden,
		Message:   autogen.Messages(message),
	}
	return r.VisitPostItemResponse(c.Response())
}
func PostItem400(c echo.Context, message string) error {
	r := autogen.PostItem400JSONResponse{
		ErrorCode: autogen.ErrBadRequest,
		Message:   autogen.Messages(message),
	}
	return r.VisitPostItemResponse(c.Response())
}

func PostItem500(c echo.Context, message string) error {
	r := autogen.PostItem500JSONResponse{
		ErrorCode: autogen.ErrInternalServerError,
		Message:   autogen.Messages(message),
	}
	return r.VisitPostItemResponse(c.Response())
}

// (POST /categories/{category_id}/items)
func (s *Server) PostItem(c echo.Context, categoryId autogen.UUID) error {
	usr, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	// Make sure the correct people have access
	switch usr.Role {
	default:
		return PostItem403(c, "You are not allowed to do this")
	case autogen.AccountMember:
	case autogen.AccountAdmin:
	case autogen.AccountSuperAdmin:
	}

	var p autogen.NewItem
	if err := c.Bind(&p); err != nil {
		logrus.Error(err)
		return PostItem400(c, "Incomprehensible body")
	}

	isMenu := false
	if p.IsMenu != nil {
		isMenu = *p.IsMenu
	}

	if !isMenu && p.MenuItems != nil && len(*p.MenuItems) > 0 {
		isMenu = true
		p.IsMenu = &isMenu
	}

	if isMenu {
		// Verify that the menu is correct
		if (p.MenuItems == nil || len(*p.MenuItems) == 0) && (p.MenuCategories == nil || len(*p.MenuCategories) == 0) {
			return PostItem400(c, "Menu must contain at least one item / category")
		}

		// Verify that all items exists
		for i, menuItem := range *p.MenuItems {
			item, err := s.DBackend.GetItem(c.Request().Context(), menuItem.Id.String())
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return PostItem400(c, "Item in menu does not exist")
				}
				logrus.Error(err)
				return PostItem500(c, "Unknown error while verifying menu")
			}
			if item.IsMenu {
				return PostItem400(c, "Menu cannot contain other menus")
			}
			menuItem.Name = item.Name
			menuItem.PictureUri = item.PictureUri
			(*p.MenuItems)[i] = menuItem
		}

		// Verify that all categories exists
		for i, menuCategory := range *p.MenuCategories {
			cat, err := s.DBackend.GetCategory(c.Request().Context(), menuCategory.Id.String())
			if err != nil {
				if err == mongo.ErrNoDocuments {
					return PostItem400(c, "Category in menu does not exist")
				}
				logrus.Error(err)
				return PostItem500(c, "Unknown error while verifying menu")
			}
			menuCategory.Name = cat.Name
			menuCategory.PictureUri = cat.PictureUri
			(*p.MenuCategories)[i] = menuCategory
		}
	}

	d, err := base64.StdEncoding.DecodeString(p.Picture)
	if err != nil {
		logrus.Error(err)
		return PostItem400(c, "Could not decode picture")
	}

	// Check MIME type
	if !strings.Contains(http.DetectContentType(d), "image") {
		logrus.Error(err)
		return PostItem400(c, "Picture is not an image")
	}

	uid := uuid.New()

	// Save image to storage
	err = storage.SaveFile("items/"+uid.String(), d)
	if err != nil {
		logrus.Error(err)
		return PostItem500(c, "Could not save picture")
	}

	item := &models.Item{
		Item: autogen.Item{
			Id:              uid,
			CategoryId:      categoryId,
			Name:            p.Name,
			Prices:          p.Prices,
			PictureUri:      "/categories/" + categoryId.String() + "/items/" + uid.String() + "/picture",
			Promotion:       p.Promotion,
			PromotionEndsAt: p.PromotionEndsAt,
			State:           p.State,
			AmountLeft:      p.AmountLeft,
			BuyLimit:        p.BuyLimit,
			OptimalAmount:   p.OptimalAmount,
			AvailableFrom:   p.AvailableFrom,
			AvailableUntil:  p.AvailableUntil,
			IsMenu:          isMenu,
			MenuItems:       p.MenuItems,
			MenuCategories:  p.MenuCategories,
		},
	}

	// Save item to database
	err = s.DBackend.CreateItem(c.Request().Context(), item)
	if err != nil {
		logrus.Error(err)
		return PostItem500(c, "Could not save item")
	}

	logrus.WithField("item", item.Name).WithField("by", usr.Name()).Info("Item created")
	autogen.PostItem201JSONResponse(item.Item).VisitPostItemResponse(c.Response())
	return nil
}

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

// (GET /categories)
func (s *Server) GetCategories(c echo.Context, p autogen.GetCategoriesParams) error {
	// Get account from cookie
	user, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	data, err := s.DBackend.GetAllCategories(c.Request().Context())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var categories []autogen.Category
	for _, category := range data {
		if category.Hidden && p.Hidden == nil {
			continue
		} else if p.Hidden != nil && *p.Hidden && !user.IsAdmin() {
			continue
		}
		categories = append(categories, category.Category)
	}

	autogen.GetCategories200JSONResponse(categories).VisitGetCategoriesResponse(c.Response())
	return nil
}

// (POST /categories)
func (s *Server) PostCategory(c echo.Context) error {
	admin, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	var p autogen.NewCategory
	if err := c.Bind(&p); err != nil {
		return Error400(c)
	}

	// Get category from request
	uid := uuid.New()

	if p.Picture != "" {
		// Get image from p.Picture as base64
		d, err := base64.StdEncoding.DecodeString(p.Picture)
		if err != nil {
			return Error400(c)
		}

		// Check MIME type
		if !strings.Contains(http.DetectContentType(d), "image") {
			return Error400(c)
		}

		// Save image to storage
		err = storage.SaveFile("categories/"+uid.String(), d)
		if err != nil {
			logrus.Error(err)
			return Error500(c)
		}
	}

	category := &models.Category{
		Category: autogen.Category{
			Id:         uid,
			Name:       p.Name,
			Position:   p.Position,
			PictureUri: "/categories/" + uid.String() + "/picture",
		},
	}

	// Save category to database
	err = s.DBackend.CreateCategory(c.Request().Context(), category)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Category %s created by %s", category.Id.String(), admin.Id.String())
	autogen.PostCategory201JSONResponse(category.Category).VisitPostCategoryResponse(c.Response())
	return nil
}

// (DELETE /categories/{category_id})
func (s *Server) MarkDeleteCategory(c echo.Context, categoryId autogen.UUID) error {
	admin, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	_, err = s.DBackend.GetCategory(c.Request().Context(), categoryId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorCategoryNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	err = s.DBackend.MarkDeleteCategory(c.Request().Context(), categoryId.String(), admin.Id.String())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Category %s marked deleted by %s", categoryId.String(), admin.Id.String())
	autogen.MarkDeleteCategory204Response{}.VisitMarkDeleteCategoryResponse(c.Response())
	return nil
}

// (GET /categories/{category_id})
func (s *Server) GetCategory(c echo.Context, categoryId autogen.UUID) error {
	// Get account from cookie
	_, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	category, err := s.DBackend.GetCategory(c.Request().Context(), categoryId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorCategoryNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	autogen.GetCategory200JSONResponse(category.Category).VisitGetCategoryResponse(c.Response())
	return nil
}

// (PATCH /categories/{category_id})
func (s *Server) PatchCategory(c echo.Context, categoryId autogen.UUID) error {
	admin, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	category, err := s.DBackend.GetCategory(c.Request().Context(), categoryId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorCategoryNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	var p autogen.UpdateCategory
	if err := c.Bind(&p); err != nil {
		return Error400(c)
	}

	if p.Name != nil {
		category.Name = *p.Name
	}

	if p.Position != nil {
		category.Position = *p.Position
	}

	if p.Hidden != nil {
		category.Hidden = *p.Hidden
	}

	if p.Picture != nil {
		// Get image from p.Picture as base64
		d, err := base64.StdEncoding.DecodeString(*p.Picture)
		if err != nil {
			return Error400(c)
		}

		// Check MIME type
		if !strings.Contains(http.DetectContentType(d), "image") {
			return Error400(c)
		}

		// Save image to storage
		err = storage.SaveFile("categories/"+category.Id.String(), d)
		if err != nil {
			logrus.Error(err)
			return Error500(c)
		}
	}

	err = s.DBackend.UpdateCategory(c.Request().Context(), category)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Category %s updated by %s", categoryId.String(), admin.Id.String())
	autogen.PatchCategory200JSONResponse(category.Category).VisitPatchCategoryResponse(c.Response())
	return nil
}

// (GET /categories/{category_id}/picture)
func (s *Server) GetCategoryPicture(c echo.Context, categoryId autogen.UUID) error {
	_, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	data, err := storage.GetFile("categories/" + categoryId.String())
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

	// Caching
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

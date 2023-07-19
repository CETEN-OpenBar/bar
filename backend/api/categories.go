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

// (GET /categories)
func (s *Server) GetCategories(c echo.Context) error {
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
		logrus.Error(err)
		return Error500(c)
	}

	data, err := s.DBackend.GetAllCategories()
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var categories []autogen.Category
	for _, category := range data {
		categories = append(categories, category.Category)
	}

	// Caching for 10 minutes
	c.Response().Header().Set("Cache-Control", "max-age=600")
	c.Response().Header().Set("Expires", "600")

	autogen.GetCategories200JSONResponse(categories).VisitGetCategoriesResponse(c.Response())
	return nil
}

// (POST /categories)
func (s *Server) PostCategory(c echo.Context) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(adminId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	var p autogen.NewCategory
	if err := c.Bind(&p); err != nil {
		return Error400(c)
	}

	// Get category from request
	uid := uuid.New()

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

	category := &models.Category{
		Category: autogen.Category{
			Id:         uid,
			Name:       p.Name,
			PictureUri: "/categories/" + uid.String() + "/picture",
		},
	}

	// Save category to database
	err = s.DBackend.CreateCategory(category)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Category %s created by %s", category.Id.String(), adminId)
	autogen.PostCategory201JSONResponse(category.Category).VisitPostCategoryResponse(c.Response())
	return nil
}

// (DELETE /categories/{category_id})
func (s *Server) MarkDeleteCategory(c echo.Context, categoryId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(adminId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}

	_, err = s.DBackend.GetCategory(categoryId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrorCategoryNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	err = s.DBackend.MarkDeleteCategory(categoryId.String(), adminId)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Category %s marked deleted by %s", categoryId.String(), adminId)
	autogen.MarkDeleteCategory204Response{}.VisitMarkDeleteCategoryResponse(c.Response())
	return nil
}

// (GET /categories/{category_id})
func (s *Server) GetCategory(c echo.Context, categoryId autogen.UUID) error {
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
		logrus.Error(err)
		return Error500(c)
	}

	category, err := s.DBackend.GetCategory(categoryId.String())
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
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminId, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get account from database
	_, err := s.DBackend.GetAccount(adminId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Delete cookie
			sess.Options.MaxAge = -1
			sess.Save(c.Request(), c.Response())
			return ErrorAccNotFound(c)
		}
		return Error500(c)
	}
	category, err := s.DBackend.GetCategory(categoryId.String())
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

	err = s.DBackend.UpdateCategory(category)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Category %s updated by %s", categoryId.String(), adminId)
	autogen.PatchCategory200JSONResponse(category.Category).VisitPatchCategoryResponse(c.Response())
	return nil
}

// (GET /categories/{category_id}/picture)
func (s *Server) GetCategoryPicture(c echo.Context, categoryId autogen.UUID) error {
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
		logrus.Error(err)
		return Error500(c)
	}

	_, err = s.DBackend.GetCategory(categoryId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Remove cache
			c.Response().Header().Set("Cache-Control", "max-age=0")
			c.Response().Header().Set("Expires", "0")
			return ErrorCategoryNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	data, err := storage.GetFile("categories/" + categoryId.String())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Caching
	c.Response().Header().Set("Cache-Control", "max-age=86400")
	c.Response().Header().Set("Expires", "86400")

	c.Response().Header().Set("Content-Type", http.DetectContentType(data))
	c.Response().Write(data)
	return nil
}

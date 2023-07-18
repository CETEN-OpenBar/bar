package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
)

// (GET /categories)
func (s *Server) GetCategories(c echo.Context) error {
	// TODO: implement
	return nil
}

// (POST /categories)
func (s *Server) PostCategory(c echo.Context) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (DELETE /categories/{category_id})
func (s *Server) DeleteCategory(c echo.Context, categoryId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (GET /categories/{category_id})
func (s *Server) GetCategory(c echo.Context, categoryId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (PATCH /categories/{category_id})
func (s *Server) PatchCategory(c echo.Context, categoryId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (GET /categories/{category_id}/picture)
func (s *Server) GetCategoryPicture(c echo.Context, categoryId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

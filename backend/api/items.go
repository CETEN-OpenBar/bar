package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
)

// (POST /categories/{category_id}/items)
func (s *Server) PostItem(c echo.Context, categoryId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (DELETE /categories/{category_id}/items/{item_id})
func (s *Server) MarkDeleteItem(c echo.Context, categoryId autogen.UUID, itemId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (PATCH /categories/{category_id}/items/{item_id})
func (s *Server) PatchItem(c echo.Context, categoryId autogen.UUID, itemId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

// (GET /categories/{category_id}/items/{item_id}/picture)
func (s *Server) GetItemPicture(c echo.Context, categoryId autogen.UUID, itemId autogen.UUID) error {
	// Get account from cookie
	sess := s.getUserSess(c)
	_, ok := sess.Values["account_id"].(string)
	if !ok {
		return Error401(c)
	}

	// TODO: implement
	return nil
}

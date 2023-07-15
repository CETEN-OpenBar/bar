package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
)

// (POST /categories/{category_id}/items)
func (s *Server) PostItem(c echo.Context, categoryId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (DELETE /categories/{category_id}/items/{item_id})
func (s *Server) DeleteItem(c echo.Context, categoryId autogen.UUID, itemId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (PATCH /categories/{category_id}/items/{item_id})
func (s *Server) PatchItem(c echo.Context, categoryId autogen.UUID, itemId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (GET /categories/{category_id}/items/{item_id}/picture)
func (s *Server) GetItemPicture(c echo.Context, categoryId autogen.UUID, itemId autogen.UUID) error {
	// TODO: implement
	return nil
}

package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
)

// (GET /carousel/images)
func (s *Server) GetCarouselImages(c echo.Context) error {
	// Get account from cookie
	sess := s.getUserSess(c)
	_, ok := sess.Values["account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (POST /carousel/images)
func (s *Server) AddCarouselImage(c echo.Context) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (DELETE /carousel/images/{image_id})
func (s *Server) MarkDeleteCarouselImage(c echo.Context, imageId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (GET /carousel/texts)
func (s *Server) GetCarouselTexts(c echo.Context) error {
	// Get account from cookie
	sess := s.getUserSess(c)
	_, ok := sess.Values["account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (POST /carousel/texts)
func (s *Server) AddCarouselText(c echo.Context) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

// (DELETE /carousel/texts/{text_id})
func (s *Server) MarkDeleteCarouselText(c echo.Context, textId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	_, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// TODO: implement
	return nil
}

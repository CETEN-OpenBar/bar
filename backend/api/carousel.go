package api

import (
	"bar/autogen"

	"github.com/labstack/echo/v4"
)

// (GET /carousel/images)
func (s *Server) GetCarouselImages(c echo.Context) error {
	// TODO: implement
	return nil
}

// (POST /carousel/images)
func (s *Server) AddCarouselImage(c echo.Context) error {
	// TODO: implement
	return nil
}

// (DELETE /carousel/images/{image_id})
func (s *Server) DeleteCarouselImage(c echo.Context, imageId autogen.UUID) error {
	// TODO: implement
	return nil
}

// (GET /carousel/texts)
func (s *Server) GetCarouselTexts(c echo.Context) error {
	// TODO: implement
	return nil
}

// (POST /carousel/texts)
func (s *Server) AddCarouselText(c echo.Context) error {
	// TODO: implement
	return nil
}

// (DELETE /carousel/texts/{text_id})
func (s *Server) DeleteCarouselText(c echo.Context, textId autogen.UUID) error {
	// TODO: implement
	return nil
}

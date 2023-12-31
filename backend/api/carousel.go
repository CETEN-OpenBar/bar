package api

import (
	"bar/autogen"
	"bar/internal/models"
	"bar/internal/storage"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// (GET /carousel/images)
func (s *Server) GetCarouselImages(c echo.Context) error {
	// Get carousel images from database
	data, err := s.DBackend.GetAllCarouselImages(c.Request().Context())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var images = make([]autogen.CarouselImage, 0)

	for _, image := range data {
		images = append(images, image.CarouselImage)
	}

	// Return carousel images
	autogen.GetCarouselImages200JSONResponse(images).VisitGetCarouselImagesResponse(c.Response())
	return nil
}

// (POST /carousel/images)
func (s *Server) AddCarouselImage(c echo.Context) error {
	admin, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	// Get image from request
	image, err := c.FormFile("image")
	if err != nil {
		return Error400(c)
	}

	file, err := image.Open()
	if err != nil {
		return Error400(c)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return Error400(c)
	}

	// Check MIME type
	if !strings.Contains(http.DetectContentType(data), "image") {
		return Error400(c)
	}

	uid := uuid.New()
	// Save image to storage
	err = storage.SaveFile("carousel/"+uid.String(), data)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	// Add image to database
	carouselImage := &models.CarouselImage{
		CarouselImage: autogen.CarouselImage{
			Id:       uid,
			ImageUrl: "/carousel/images/" + uid.String(),
		},
	}

	err = s.DBackend.CreateCarouselImage(c.Request().Context(), carouselImage)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("image", carouselImage.Id.String()).WithField("by", admin.Name()).Info("Carousel image added")
	autogen.AddCarouselImage201JSONResponse(carouselImage.CarouselImage).VisitAddCarouselImageResponse(c.Response())
	return nil
}

// (GET /carousel/images/{image_id})
func (s *Server) GetCarouselImage(c echo.Context, imageId autogen.UUID) error {
	data, err := storage.GetFile("carousel/" + imageId.String())
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

// (DELETE /carousel/images/{image_id})
func (s *Server) MarkDeleteCarouselImage(c echo.Context, imageId autogen.UUID) error {
	admin, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	_, err = s.DBackend.GetCarouselImage(c.Request().Context(), imageId.String())
	if err != nil {
		return ErrorImageNotFound(c)
	}

	err = s.DBackend.MarkDeleteCarouselImage(c.Request().Context(), imageId.String(), admin.Id.String())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("image", imageId.String()).WithField("by", admin.Name()).Info("Carousel image marked for deletion")
	autogen.MarkDeleteAccountId204Response{}.VisitMarkDeleteAccountIdResponse(c.Response())
	return nil
}

// (GET /carousel/texts)
func (s *Server) GetCarouselTexts(c echo.Context) error {
	// Get carousel images from database
	data, err := s.DBackend.GetAllCarouselTexts(c.Request().Context())
	if err != nil {
		return Error500(c)
	}

	var texts = make([]autogen.CarouselText, 0)

	for _, image := range data {
		texts = append(texts, image.CarouselText)
	}

	// Return carousel images
	autogen.GetCarouselTexts200JSONResponse(texts).VisitGetCarouselTextsResponse(c.Response())
	return nil
}

// (POST /carousel/texts)
func (s *Server) AddCarouselText(c echo.Context) error {
	admin, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	// Get text from request
	var text autogen.CarouselTextCreate
	err = c.Bind(&text)
	if err != nil {
		return Error400(c)
	}

	var color = "#000000"
	if text.Color != nil {
		color = *text.Color
	}

	t := &models.CarouselText{
		CarouselText: autogen.CarouselText{
			Id:    uuid.New(),
			Text:  text.Text,
			Color: color,
		},
	}

	err = s.DBackend.CreateCarouselText(c.Request().Context(), t)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("text", t.Id.String()).WithField("by", admin.Name()).Info("Carousel text added")
	autogen.AddCarouselText201JSONResponse(t.CarouselText).VisitAddCarouselTextResponse(c.Response())
	return nil
}

// (DELETE /carousel/texts/{text_id})
func (s *Server) MarkDeleteCarouselText(c echo.Context, textId autogen.UUID) error {
	admin, err := MustGetAdmin(c)
	if err != nil {
		return nil
	}

	_, err = s.DBackend.GetCarouselText(c.Request().Context(), textId.String())
	if err != nil {
		return ErrorTextNotFound(c)
	}

	err = s.DBackend.MarkDeleteCarouselText(c.Request().Context(), textId.String(), admin.Id.String())
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.WithField("text", textId.String()).WithField("by", admin.Name()).Info("Carousel text marked for deletion")
	return nil
}

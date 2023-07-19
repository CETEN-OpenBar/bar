package api

import (
	"bar/autogen"
	"bar/internal/models"
	"bar/internal/storage"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// (GET /carousel/images)
func (s *Server) GetCarouselImages(c echo.Context) error {
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

	// Get carousel images from database
	data, err := s.DBackend.GetAllCarouselImages()
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	var images []autogen.CarouselImage

	for _, image := range data {
		images = append(images, image.CarouselImage)
	}

	// Caching for 10 minutes
	c.Response().Header().Set("Cache-Control", "max-age=600")
	c.Response().Header().Set("Expires", "600")

	// Return carousel images
	autogen.GetCarouselImages200JSONResponse(images).VisitGetCarouselImagesResponse(c.Response())
	return nil
}

// (POST /carousel/images)
func (s *Server) AddCarouselImage(c echo.Context) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminID, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
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

	err = s.DBackend.CreateCarouselImage(carouselImage)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Carousel image %s added by admin %s", uid.String(), adminID)
	autogen.AddCarouselImage201JSONResponse(carouselImage.CarouselImage).VisitAddCarouselImageResponse(c.Response())
	return nil
}

// (GET /carousel/images/{image_id})
func (s *Server) GetCarouselImage(c echo.Context, imageId autogen.UUID) error {
	_, err := s.DBackend.GetCarouselImage(imageId.String())
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Remove cache
			c.Response().Header().Set("Cache-Control", "max-age=0")
			c.Response().Header().Set("Expires", "0")
			return ErrorImageNotFound(c)
		}
		logrus.Error(err)
		return Error500(c)
	}

	data, err := storage.GetFile("carousel/" + imageId.String())
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

// (DELETE /carousel/images/{image_id})
func (s *Server) MarkDeleteCarouselImage(c echo.Context, imageId autogen.UUID) error {
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

	_, err = s.DBackend.GetCarouselImage(imageId.String())
	if err != nil {
		return ErrorImageNotFound(c)
	}

	err = s.DBackend.MarkDeleteCarouselImage(imageId.String(), adminId)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Carousel image %s marked for deletion by admin %s", imageId.String(), adminId)
	autogen.MarkDeleteAccountId204Response{}.VisitMarkDeleteAccountIdResponse(c.Response())
	return nil
}

// (GET /carousel/texts)
func (s *Server) GetCarouselTexts(c echo.Context) error {
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

	// Get carousel images from database
	data, err := s.DBackend.GetAllCarouselTexts()
	if err != nil {
		return Error500(c)
	}

	var texts []autogen.CarouselText

	for _, image := range data {
		texts = append(texts, image.CarouselText)
	}

	// Caching for 10 minutes
	c.Response().Header().Set("Cache-Control", "max-age=600")
	c.Response().Header().Set("Expires", "600")

	// Return carousel images
	autogen.GetCarouselTexts200JSONResponse(texts).VisitGetCarouselTextsResponse(c.Response())
	return nil
}

// (POST /carousel/texts)
func (s *Server) AddCarouselText(c echo.Context) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminID, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	// Get text from request
	var text autogen.CarouselTextCreate
	err := c.Bind(&text)
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

	err = s.DBackend.CreateCarouselText(t)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Carousel text %s added by admin %s", t.Id.String(), adminID)
	autogen.AddCarouselText201JSONResponse(t.CarouselText).VisitAddCarouselTextResponse(c.Response())
	return nil
}

// (DELETE /carousel/texts/{text_id})
func (s *Server) MarkDeleteCarouselText(c echo.Context, textId autogen.UUID) error {
	// Get admin account from cookie
	sess := s.getAdminSess(c)
	adminID, ok := sess.Values["admin_account_id"].(string)
	if !ok {
		return ErrorNotAuthenticated(c)
	}

	_, err := s.DBackend.GetCarouselText(textId.String())
	if err != nil {
		return ErrorTextNotFound(c)
	}

	err = s.DBackend.MarkDeleteCarouselText(textId.String(), adminID)
	if err != nil {
		logrus.Error(err)
		return Error500(c)
	}

	logrus.Infof("Carousel text %s marked for deletion by admin %s", textId.String(), adminID)
	return nil
}

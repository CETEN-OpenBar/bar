package api

import (
	"bar/autogen"
	"context"
)

// (GET /carousel/images)
func (s *Server) GetCarouselImages(ctx context.Context, request autogen.GetCarouselImagesRequestObject) (autogen.GetCarouselImagesResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (POST /carousel/images)
func (s *Server) AddCarouselImage(ctx context.Context, request autogen.AddCarouselImageRequestObject) (autogen.AddCarouselImageResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (DELETE /carousel/images/{image_id})
func (s *Server) DeleteCarouselImage(ctx context.Context, request autogen.DeleteCarouselImageRequestObject) (autogen.DeleteCarouselImageResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /carousel/texts)
func (s *Server) GetCarouselTexts(ctx context.Context, request autogen.GetCarouselTextsRequestObject) (autogen.GetCarouselTextsResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (POST /carousel/texts)
func (s *Server) AddCarouselText(ctx context.Context, request autogen.AddCarouselTextRequestObject) (autogen.AddCarouselTextResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (DELETE /carousel/texts/{text_id})
func (s *Server) DeleteCarouselText(ctx context.Context, request autogen.DeleteCarouselTextRequestObject) (autogen.DeleteCarouselTextResponseObject, error) {
	// TODO: implement
	return nil, nil
}

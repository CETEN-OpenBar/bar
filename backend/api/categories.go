package api

import (
	"bar/autogen"
	"context"
)

// (GET /categories)
func (s *Server) GetCategories(ctx context.Context, request autogen.GetCategoriesRequestObject) (autogen.GetCategoriesResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (POST /categories)
func (s *Server) PostCategory(ctx context.Context, request autogen.PostCategoryRequestObject) (autogen.PostCategoryResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (DELETE /categories/{category_id})
func (s *Server) DeleteCategory(ctx context.Context, request autogen.DeleteCategoryRequestObject) (autogen.DeleteCategoryResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /categories/{category_id})
func (s *Server) GetCategory(ctx context.Context, request autogen.GetCategoryRequestObject) (autogen.GetCategoryResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (PATCH /categories/{category_id})
func (s *Server) PatchCategory(ctx context.Context, request autogen.PatchCategoryRequestObject) (autogen.PatchCategoryResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /categories/{category_id}/picture)
func (s *Server) GetCategoryPicture(ctx context.Context, request autogen.GetCategoryPictureRequestObject) (autogen.GetCategoryPictureResponseObject, error) {
	// TODO: implement
	return nil, nil
}

package api

import (
	"bar/autogen"
	"context"
)

// (POST /categories/{category_id}/items)
func (s *Server) PostItem(ctx context.Context, request autogen.PostItemRequestObject) (autogen.PostItemResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (DELETE /categories/{category_id}/items/{item_id})
func (s *Server) DeleteItem(ctx context.Context, request autogen.DeleteItemRequestObject) (autogen.DeleteItemResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (PATCH /categories/{category_id}/items/{item_id})
func (s *Server) PatchItem(ctx context.Context, request autogen.PatchItemRequestObject) (autogen.PatchItemResponseObject, error) {
	// TODO: implement
	return nil, nil
}

// (GET /categories/{category_id}/items/{item_id}/picture)
func (s *Server) GetItemPicture(ctx context.Context, request autogen.GetItemPictureRequestObject) (autogen.GetItemPictureResponseObject, error) {
	// TODO: implement
	return nil, nil
}

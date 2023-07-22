package mongo

import (
	"bar/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (b *Backend) GetAllCarouselImages(ctx context.Context) ([]*models.CarouselImage, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var carouselImages []*models.CarouselImage
	cursor, err := b.db.Collection(CarouselImagesCollection).Find(ctx, bson.M{
		"$or": []bson.M{
			{
				"deleted_at": bson.M{
					"$exists": false,
				},
			},
			{
				"deleted_at": nil,
			},
		},
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return carouselImages, nil
		}
		return nil, err
	}

	if err := cursor.All(ctx, &carouselImages); err != nil {
		return nil, err
	}

	return carouselImages, nil
}

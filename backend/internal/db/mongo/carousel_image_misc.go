package mongo

import (
	"bar/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (b *Backend) GetAllCarouselImages() ([]*models.CarouselImage, error) {
	ctx, cancel := b.GetContext()
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

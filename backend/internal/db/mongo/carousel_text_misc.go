package mongo

import (
	"bar/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (b *Backend) GetAllCarouselTexts(ctx context.Context) ([]*models.CarouselText, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var carouselTexts []*models.CarouselText
	cursor, err := b.db.Collection(CarouselTextsCollection).Find(ctx, bson.M{
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
			return carouselTexts, nil
		}
		return nil, err
	}

	if err := cursor.All(ctx, &carouselTexts); err != nil {
		return nil, err
	}

	return carouselTexts, nil
}

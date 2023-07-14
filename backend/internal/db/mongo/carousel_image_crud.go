package mongo

import (
	"bar/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) CreateCarouselImage(ci *models.CarouselImage) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(CarouselImagesCollection).InsertOne(ctx, ci)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) GetCarouselImage(id string) (*models.CarouselImage, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var ci models.CarouselImage
	err := b.db.Collection(CarouselImagesCollection).FindOne(ctx,
		bson.M{
			"id": id,
		},
	).Decode(&ci)
	if err != nil {
		return nil, err
	}

	return &ci, nil
}

func (b *Backend) UpdateCarouselImage(ci *models.CarouselImage) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselImagesCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": ci.Id,
		},
		bson.M{
			"$set": ci,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) DeleteCarouselImage(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(CarouselImagesCollection).DeleteOne(ctx,
		bson.M{
			"id": id,
		})
	if err != nil {
		return err
	}

	return nil
}

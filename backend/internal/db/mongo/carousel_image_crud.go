package mongo

import (
	"bar/internal/models"
	"time"

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
			"deleted_at": bson.M{
				"$exists": false,
			},
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
			"deleted_at": bson.M{
				"$exists": false,
			},
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

func (b *Backend) DeleteCarouselImage(id, by string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselImagesCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": id,
		},
		bson.M{
			"$set": bson.M{
				"deleted_at": time.Now().Unix(),
				"deleted_by": by,
			},
		},
		options.FindOneAndUpdate().SetUpsert(false))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

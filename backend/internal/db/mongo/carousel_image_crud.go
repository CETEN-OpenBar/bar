package mongo

import (
	"bar/internal/models"
	"time"

	"github.com/google/uuid"
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
			"id": uuid.MustParse(id),

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

func (b *Backend) MarkDeleteCarouselImage(id, by string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselImagesCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
		bson.M{
			"$set": bson.M{
				"deleted_at": time.Now().Unix(),
				"deleted_by": uuid.MustParse(id),
			},
		},
		options.FindOneAndUpdate().SetUpsert(false))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) DeleteCarouselImage(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselImagesCollection).FindOneAndDelete(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) RestoreCarouselImage(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselImagesCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
		bson.M{
			"$unset": bson.M{
				"deleted_at": "",
				"deleted_by": "",
			},
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

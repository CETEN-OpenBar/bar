package mongo

import (
	"bar/internal/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) CreateCarouselText(ct *models.CarouselText) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(CarouselTextsCollection).InsertOne(ctx, ct)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) GetCarouselText(id string) (*models.CarouselText, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var ct models.CarouselText
	err := b.db.Collection(CarouselTextsCollection).FindOne(ctx,
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
	).Decode(&ct)
	if err != nil {
		return nil, err
	}

	return &ct, nil
}

func (b *Backend) UpdateCarouselText(ct *models.CarouselText) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselTextsCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": ct.Id,

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
			"$set": ct,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) MarkDeleteCarouselText(id, by string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselTextsCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": uuid.MustParse(id),
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

func (b *Backend) DeleteCarouselText(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselTextsCollection).FindOneAndDelete(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) RestoreCarouselText(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CarouselTextsCollection).FindOneAndUpdate(ctx,
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

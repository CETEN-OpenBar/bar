package mongo

import (
	"bar/internal/models"

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
			"id": id,
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

func (b *Backend) DeleteCarouselText(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(CarouselTextsCollection).DeleteOne(ctx,
		bson.M{
			"id": id,
		})
	if err != nil {
		return err
	}

	return nil
}

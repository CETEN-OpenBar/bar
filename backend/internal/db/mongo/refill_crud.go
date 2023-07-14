package mongo

import (
	"bar/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) CreateRefill(refill *models.Refill) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(RefillsCollection).InsertOne(ctx, refill)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) GetRefill(id string) (*models.Refill, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var refill models.Refill
	err := b.db.Collection(RefillsCollection).FindOne(ctx,
		bson.M{
			"id": id,
		},
	).Decode(&refill)
	if err != nil {
		return nil, err
	}

	return &refill, nil
}

func (b *Backend) UpdateRefill(refill *models.Refill) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(RefillsCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": refill.Id,
		},
		bson.M{
			"$set": refill,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) DeleteRefill(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(RefillsCollection).DeleteOne(ctx,
		bson.M{
			"id": id,
		})
	if err != nil {
		return err
	}

	return nil
}

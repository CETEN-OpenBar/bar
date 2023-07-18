package mongo

import (
	"bar/internal/models"
	"time"

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
			"deleted_at": bson.M{
				"$exists": false,
			},
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
			"deleted_at": bson.M{
				"$exists": false,
			},
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

func (b *Backend) MarkDeleteRefill(id, by string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(RefillsCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) DeleteRefill(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(RefillsCollection).FindOneAndDelete(ctx,
		bson.M{
			"id": id,
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) RestoreRefill(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(RefillsCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": id,
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

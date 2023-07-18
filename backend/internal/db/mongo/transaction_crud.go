package mongo

import (
	"bar/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) CreateTransaction(tx *models.Transaction) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(TransactionsCollection).InsertOne(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) GetTransaction(id string) (*models.Transaction, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var tx models.Transaction
	err := b.db.Collection(TransactionsCollection).FindOne(ctx,
		bson.M{
			"id": id,
			"deleted_at": bson.M{
				"$exists": false,
			},
		},
	).Decode(&tx)
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (b *Backend) UpdateTransaction(tx *models.Transaction) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(TransactionsCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": tx.Id,
			"deleted_at": bson.M{
				"$exists": false,
			},
		},
		bson.M{
			"$set": tx,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) MarkDeleteTransaction(id, by string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(TransactionsCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) DeleteTransaction(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(TransactionsCollection).FindOneAndDelete(ctx,
		bson.M{
			"id": id,
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) RestoreTransaction(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(TransactionsCollection).FindOneAndUpdate(ctx,
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

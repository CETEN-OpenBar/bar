package mongo

import (
	"bar/internal/models"

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

func (b *Backend) DeleteTransaction(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(TransactionsCollection).DeleteOne(ctx,
		bson.M{
			"id": id,
		})
	if err != nil {
		return err
	}

	return nil
}

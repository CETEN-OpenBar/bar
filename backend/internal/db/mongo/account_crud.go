package mongo

import (
	"bar/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) CreateAccount(acc *models.Account) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(AccountsCollection).InsertOne(ctx, acc)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) GetAccount(id string) (*models.Account, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var acc models.Account
	err := b.db.Collection(AccountsCollection).FindOne(ctx,
		bson.M{
			"id": id,
		},
	).Decode(&acc)
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

func (b *Backend) UpdateAccount(acc *models.Account) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(AccountsCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": acc.Id,
		},
		bson.M{
			"$set": acc,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) DeleteAccount(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(AccountsCollection).DeleteOne(ctx,
		bson.M{
			"id": id,
		})
	if err != nil {
		return err
	}

	return nil
}

package mongo

import (
	"bar/internal/models"
	"time"

	"github.com/google/uuid"
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
			"$set": acc,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) MarkDeleteAccount(id, by string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	// Mark deleted_at
	res := b.db.Collection(AccountsCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) DeleteAccount(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(AccountsCollection).FindOneAndDelete(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) RestoreAccount(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(AccountsCollection).FindOneAndUpdate(ctx,
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

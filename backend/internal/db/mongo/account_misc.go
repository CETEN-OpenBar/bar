package mongo

import (
	"bar/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func (b *Backend) GetAccounts(page int, size int) ([]*models.Account, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	// Get "size" accounts from "page" using aggregation
	var accounts []*models.Account
	cursor, err := b.db.Collection(AccountsCollection).Aggregate(ctx, []bson.M{
		{
			"$skip": page * size,
		},
		{
			"$limit": size,
		},
	})
	if err != nil {
		return nil, err
	}

	// Decode each account
	if err := cursor.All(ctx, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (b *Backend) CountAccounts() (int64, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	// Count all accounts
	return b.db.Collection(AccountsCollection).CountDocuments(ctx, bson.M{})
}

package mongo

import (
	"bar/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func (b *Backend) GetAccounts(page uint64, size uint64) ([]*models.Account, error) {
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

func (b *Backend) CountAccounts() (uint64, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	// Count all accounts
	count, err := b.db.Collection(AccountsCollection).CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (b *Backend) GetAccountByCard(card string) (*models.Account, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	// Get account by card
	var account models.Account
	if err := b.db.Collection(AccountsCollection).FindOne(ctx, bson.M{"card_id": card}).Decode(&account); err != nil {
		return nil, err
	}

	return &account, nil
}
func (b *Backend) GetAccountByGoogle(googleID string) (*models.Account, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	// Get account by card
	var account models.Account
	if err := b.db.Collection(AccountsCollection).FindOne(ctx, bson.M{"google_id": googleID}).Decode(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

package mongo

import (
	"bar/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func (b *Backend) GetTransactions(accountID string, page int, size int, state string) ([]*models.Transaction, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	// Get "size" transactions from "page" using aggregation
	var transactions []*models.Transaction
	cursor, err := b.db.Collection(TransactionsCollection).Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"account_id": accountID,
				"state":      state,
			},
		},
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
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

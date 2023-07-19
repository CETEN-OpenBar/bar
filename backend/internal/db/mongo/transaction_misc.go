package mongo

import (
	"bar/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) GetTransactions(accountID string, page uint64, size uint64, state string) ([]*models.Transaction, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	// Get "size" transactions from "page" using aggregation
	var transactions []*models.Transaction
	cursor, err := b.db.Collection(TransactionsCollection).Find(ctx, bson.M{
		"account_id": accountID,
		"state":      state,
	}, options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)))
	if err != nil {
		return nil, err
	}

	// Decode each account
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (b *Backend) CountTransactions(accountID string, state string) (uint64, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	count, err := b.db.Collection(TransactionsCollection).CountDocuments(ctx, bson.M{
		"account_id": accountID,
		"state":      state,
	})
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (b *Backend) GetAllTransactions(page uint64, size uint64, state string) ([]*models.Transaction, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	// Get "size" transactions from "page" using aggregation
	var transactions []*models.Transaction
	cursor, err := b.db.Collection(TransactionsCollection).Find(ctx, bson.M{
		"state": state,
	}, options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)))
	if err != nil {
		return nil, err
	}

	// Decode each account
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (b *Backend) CountAllTransactions(state string) (uint64, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	count, err := b.db.Collection(TransactionsCollection).CountDocuments(ctx, bson.M{
		"state": state,
	})
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

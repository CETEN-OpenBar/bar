package mongo

import (
	"bar/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) GetTransactions(ctx context.Context, accountID string, page uint64, size uint64, state string) ([]*models.Transaction, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{
		"account_id": accountID,
	}

	if state != "" {
		filter["state"] = state
	}

	// Get "size" transactions from "page" using aggregation
	var transactions []*models.Transaction
	cursor, err := b.db.Collection(TransactionsCollection).Find(ctx, filter, options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}

	// Decode each account
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (b *Backend) CountTransactions(ctx context.Context, accountID string, state string) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{
		"account_id": accountID,
	}

	if state != "" {
		filter["state"] = state
	}

	count, err := b.db.Collection(TransactionsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (b *Backend) GetAllTransactions(ctx context.Context, page uint64, size uint64, state string) ([]*models.Transaction, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{}

	if state != "" {
		filter["state"] = state
	}

	// Get "size" transactions from "page" using aggregation
	var transactions []*models.Transaction
	cursor, err := b.db.Collection(TransactionsCollection).Find(ctx, filter, options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}

	// Decode each account
	if err := cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (b *Backend) CountAllTransactions(ctx context.Context, state string) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{}

	if state != "" {
		filter["state"] = state
	}

	count, err := b.db.Collection(TransactionsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

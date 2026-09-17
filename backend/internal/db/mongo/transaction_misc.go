package mongo

import (
	"bar/autogen"
	"bar/internal/models"
	"context"
	"regexp"

	"github.com/google/uuid"

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

func (b *Backend) GetAllTransactions(ctx context.Context, page uint64, size uint64, TransactionState string, hide_canceled bool, name string, hide_remotes bool, StartTime int, EndTime int, ItemID string) ([]*models.Transaction, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := allTransactionsFilter(TransactionState, hide_canceled, name, hide_remotes, StartTime, EndTime, ItemID)

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

func (b *Backend) CountAllTransactions(ctx context.Context, TransactionState string, hide_canceled bool, name string, hide_remotes bool, StartTime int, EndTime int, ItemID string) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := allTransactionsFilter(TransactionState, hide_canceled, name, hide_remotes, StartTime, EndTime, ItemID)

	count, err := b.db.Collection(TransactionsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

// Keep the result list and its pagination count scoped to the same transactions.
func allTransactionsFilter(state string, hideCanceled bool, name string, hideRemotes bool, startTime int, endTime int, itemID string) bson.M {
	filter := bson.M{}

	if state != "" {
		filter["state"] = state
	}

	if name != "" {
		pattern := regexp.QuoteMeta(name)
		filter["$or"] = []bson.M{
			{
				"account_name": bson.M{
					"$regex":   pattern,
					"$options": "i",
				},
			},
			{
				"account_nick_name": bson.M{
					"$regex":   pattern,
					"$options": "i",
				},
			},
		}
	}

	if hideRemotes {
		// $ne also includes legacy transactions with a null or missing field.
		filter["is_remote"] = bson.M{"$ne": true}
	}

	if startTime > 0 || endTime > 0 {
		timeFilter := bson.M{}
		if startTime > 0 {
			timeFilter["$gte"] = startTime
		}
		if endTime > 0 {
			timeFilter["$lte"] = endTime
		}
		filter["created_at"] = timeFilter
	}

	itemsElemMatch := bson.M{}
	if itemID != "" {
		itemsElemMatch["item_id"] = uuid.MustParse(itemID)
	}
	if hideCanceled {
		itemsElemMatch["state"] = bson.M{"$ne": autogen.TransactionItemCanceled}
	}
	if len(itemsElemMatch) > 0 {
		filter["items"] = bson.M{"$elemMatch": itemsElemMatch}
	}

	return filter
}

func (b *Backend) GetAllActiveTransactionsItems(ctx context.Context, name string) ([]autogen.TransactionItem, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{"state": autogen.TransactionStarted, "items.state": autogen.TransactionItemStarted}

	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	cursor, err := b.db.Collection(TransactionsCollection).Aggregate(ctx, []bson.M{
		{"$match": filter},
		{
			"$project": bson.M{
				"_id":   0,
				"items": 1,
			},
		},
		{"$unwind": "$items"},
		{
			"$group": bson.M{
				"_id":  "$items.item_id",
				"item": bson.M{"$first": "$items"},
				"total_amount": bson.M{
					"$sum": "$items.item_amount",
				},
				"already_done": bson.M{
					"$sum": "$items.item_already_done",
				},
			},
		},
		{
			"$set": bson.M{
				"item.item_amount":       "$total_amount",
				"item.item_already_done": "$already_done",
			},
		},
		{"$replaceRoot": bson.M{"newRoot": "$item"}},
		{"$sort": bson.M{"item_amount": -1}},
	})
	if err != nil {
		return nil, err
	}

	var items []autogen.TransactionItem

	// Decode each item
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

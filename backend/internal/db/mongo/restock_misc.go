package mongo

import (
	"bar/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) GetRestocks(ctx context.Context, accountID string, page uint64, size uint64) ([]*models.Restock, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{
		"account_id": accountID,
	}

	// Get "size" restocks from "page" using aggregation
	var restocks []*models.Restock
	cursor, err := b.db.Collection(RestocksCollection).Find(ctx, filter, options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}

	// Decode each account
	if err := cursor.All(ctx, &restocks); err != nil {
		return nil, err
	}

	return restocks, nil
}

func (b *Backend) CountRestocks(ctx context.Context, accountID string) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{
		"account_id": accountID,
	}

	count, err := b.db.Collection(RestocksCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (b *Backend) GetAllRestocks(ctx context.Context, page uint64, size uint64) ([]*models.Restock, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{
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
	}

	// Get "size" restocks from "page" using aggregation
	var restocks []*models.Restock
	cursor, err := b.db.Collection(RestocksCollection).Find(ctx, filter, options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}

	// Decode each account
	if err := cursor.All(ctx, &restocks); err != nil {
		return nil, err
	}

	return restocks, nil
}

func (b *Backend) CountAllRestocks(ctx context.Context) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{
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
	}

	count, err := b.db.Collection(RestocksCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

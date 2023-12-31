package mongo

import (
	"bar/internal/models"
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) GetRefills(ctx context.Context, account string, page uint64, size uint64, startAt, endAt uint64) ([]*models.Refill, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	// Get "size" refills from "page" using aggregation
	var refills []*models.Refill
	cursor, err := b.db.Collection(RefillsCollection).Find(ctx, bson.M{
		"account_id": uuid.MustParse(account),
		"issued_at": bson.M{
			"$gte": startAt,
			"$lte": endAt,
		},
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
	}, options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &refills); err != nil {
		return nil, err
	}

	return refills, nil
}

func (b *Backend) CountRefills(ctx context.Context, account string, startAt, endAt uint64) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	count, err := b.db.Collection(RefillsCollection).CountDocuments(ctx, bson.M{
		"account_id": uuid.MustParse(account),
		"issued_at": bson.M{
			"$gte": startAt,
			"$lte": endAt,
		},
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
	})
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (b *Backend) GetAllRefills(ctx context.Context, page uint64, size uint64, startAt, endAt uint64) ([]*models.Refill, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	// Get "size" refills from "page" using aggregation
	var refills []*models.Refill
	cursor, err := b.db.Collection(RefillsCollection).Find(ctx, bson.M{
		"issued_at": bson.M{
			"$gte": startAt,
			"$lte": endAt,
		},
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
	}, options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &refills); err != nil {
		return nil, err
	}

	return refills, nil
}

func (b *Backend) CountAllRefills(ctx context.Context, startAt, endAt uint64) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	count, err := b.db.Collection(RefillsCollection).CountDocuments(ctx, bson.M{
		"issued_at": bson.M{
			"$gte": startAt,
			"$lte": endAt,
		},
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
	})
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

package mongo

import (
	"bar/internal/models"
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) GetStarrings(ctx context.Context, account string, page uint64, size uint64, startAt, endAt uint64) ([]*models.Starring, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var starrings []*models.Starring
	cursor, err := b.db.Collection(StarringCollection).Find(ctx, bson.M{
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

	if err := cursor.All(ctx, &starrings); err != nil {
		return nil, err
	}

	return starrings, nil
}

func (b *Backend) CountStarrings(ctx context.Context, account string, startAt, endAt uint64) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	count, err := b.db.Collection(StarringCollection).CountDocuments(ctx, bson.M{
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

func (b *Backend) GetAllStarrings(ctx context.Context, page uint64, size uint64, name string, startAt, endAt uint64) ([]*models.Starring, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var starrings []*models.Starring
	filter := bson.M{
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
	}

	if name != "" {
		filter["account_name"] = bson.M{"$regex": name, "$options": "i"}
	}

	cursor, err := b.db.Collection(StarringCollection).Find(ctx, filter,
		options.Find().
			SetSkip(int64(page*size)).
			SetLimit(int64(size)).
			SetSort(bson.M{"created_at": -1}),
	)

	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &starrings); err != nil {
		return nil, err
	}

	return starrings, nil
}

func (b *Backend) CountAllStarrings(ctx context.Context, name string, startAt, endAt uint64) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{
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
	}

	if name != "" {
		filter["account_name"] = bson.M{"$regex": name, "$options": "i"}
	}

	count, err := b.db.Collection(StarringCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

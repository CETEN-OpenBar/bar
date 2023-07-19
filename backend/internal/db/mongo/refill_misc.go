package mongo

import (
	"bar/internal/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (b *Backend) GetRefills(account string, page int, size int, startAt, endAt int64) ([]*models.Refill, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	// Get "size" refills from "page" using aggregation
	var refills []*models.Refill
	cursor, err := b.db.Collection(RefillsCollection).Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
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

	if err := cursor.All(ctx, &refills); err != nil {
		return nil, err
	}

	return refills, nil
}

func (b *Backend) CountRefills(account string, startAt, endAt int64) (int64, error) {
	ctx, cancel := b.GetContext()
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

	return count, nil
}

func (b *Backend) GetAllRefills(page int, size int, startAt, endAt int64) ([]*models.Refill, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	// Get "size" refills from "page" using aggregation
	var refills []*models.Refill
	cursor, err := b.db.Collection(RefillsCollection).Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
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

	if err := cursor.All(ctx, &refills); err != nil {
		return nil, err
	}

	return refills, nil
}

func (b *Backend) CountAllRefills(startAt, endAt int64) (int64, error) {
	ctx, cancel := b.GetContext()
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

	return count, nil
}

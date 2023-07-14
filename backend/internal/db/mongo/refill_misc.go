package mongo

import (
	"bar/internal/models"

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
				"account_id": account,
				"issued_at": bson.M{
					"$gte": startAt,
					"$lte": endAt,
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

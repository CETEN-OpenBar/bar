package mongo

import (
	"bar/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) GetLatestCashMovement(ctx context.Context) (*models.CashMovement, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	// Get "size" cashmovements from "page" using aggregation
	var cashmovements []*models.CashMovement
	cursor, err := b.db.Collection(CashMovementsCollection).Find(ctx, bson.M{}, options.Find().SetSkip(0).SetLimit(1).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}

	// Decode each account
	if err := cursor.All(ctx, &cashmovements); err != nil {
		return nil, err
	}

	if len(cashmovements) == 0 {
		return &models.CashMovement{}, nil
	}

	return cashmovements[0], nil
}

func (b *Backend) GetAllCashMovements(ctx context.Context, page uint64, size uint64, search string) ([]*models.CashMovement, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{}

	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"created_by_name": bson.M{"$regex": search, "$options": "i"}},
				{"reason": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	// Get "size" cashmovements from "page" using aggregation
	var cashmovements []*models.CashMovement
	cursor, err := b.db.Collection(CashMovementsCollection).Find(ctx, filter, options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return nil, err
	}

	// Decode each account
	if err := cursor.All(ctx, &cashmovements); err != nil {
		return nil, err
	}

	return cashmovements, nil
}

func (b *Backend) CountAllCashMovements(ctx context.Context, search string) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{}

	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"created_by_name": bson.M{"$regex": search, "$options": "i"}},
				{"reason": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	count, err := b.db.Collection(CashMovementsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

package mongo

import (
	"bar/internal/models"
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) GetItems(ctx context.Context, categoryID string, page, size uint64, state string) ([]*models.Item, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var items []*models.Item

	filter := bson.M{
		"category_id": uuid.MustParse(categoryID),
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
	if state != "" {
		filter["state"] = state
	}

	cursor, err := b.db.Collection(ItemsCollection).Find(ctx, filter, options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)))
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (b *Backend) CountItems(ctx context.Context, categoryID string, state string) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	count, err := b.db.Collection(ItemsCollection).CountDocuments(ctx, bson.M{
		"category_id": uuid.MustParse(categoryID),
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
		"state": state,
	})
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

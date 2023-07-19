package mongo

import (
	"bar/internal/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func (b *Backend) GetItems(categoryID string, page, size uint64, state string) ([]*models.Item, error) {
	ctx, cancel := b.GetContext()
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

	cursor, err := b.db.Collection(ItemsCollection).Aggregate(ctx, []bson.M{
		{
			"$match": filter,
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

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (b *Backend) CountItems(categoryID string, state string) (int64, error) {
	ctx, cancel := b.GetContext()
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

	return count, nil
}

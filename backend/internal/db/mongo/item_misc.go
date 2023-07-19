package mongo

import (
	"bar/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func (b *Backend) GetAllItems(categoryID string) ([]*models.Item, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var items []*models.Item
	cursor, err := b.db.Collection(ItemsCollection).Find(ctx, bson.M{
		"category_id": categoryID,
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
		return nil, err
	}

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

package mongo

import (
	"bar/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (b *Backend) GetAllCategories() ([]*models.Category, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var categories []*models.Category
	cursor, err := b.db.Collection(CategoriesCollection).Find(ctx, bson.M{
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
		if err == mongo.ErrNoDocuments {
			return categories, nil
		}
		return nil, err
	}

	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

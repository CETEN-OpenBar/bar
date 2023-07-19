package mongo

import (
	"bar/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (b *Backend) GetDeletedCategories(page uint64, size uint64) ([]*models.Category, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var accs []*models.Category
	cursor, err := b.db.Collection(CategoriesCollection).Find(ctx,
		bson.M{
			"deleted_at": bson.M{
				"$ne": nil,
			},
		},
		options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)))
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &accs); err != nil {
		return nil, err
	}

	return accs, nil
}

func (b *Backend) CountDeletedCategories() (int64, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	count, err := b.db.Collection(CategoriesCollection).CountDocuments(ctx, bson.M{
		"deleted_at": bson.M{
			"$ne": nil,
		},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

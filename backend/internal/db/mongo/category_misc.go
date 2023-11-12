package mongo

import (
	"bar/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var categories []*models.Category
	// All categories sorted by position
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
	}, options.Find().SetSort(bson.M{"position": 1}))
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

func (b *Backend) GetDeletedCategories(ctx context.Context, page uint64, size uint64) ([]*models.Category, error) {
	ctx, cancel := b.TimeoutContext(ctx)
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

func (b *Backend) CountDeletedCategories(ctx context.Context) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	count, err := b.db.Collection(CategoriesCollection).CountDocuments(ctx, bson.M{
		"deleted_at": bson.M{
			"$ne": nil,
		},
	})
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

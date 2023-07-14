package mongo

import (
	"bar/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) CreateCategory(c *models.Category) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(CategoriesCollection).InsertOne(ctx, c)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) GetCategory(id string) (*models.Category, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var c models.Category
	err := b.db.Collection(CategoriesCollection).FindOne(ctx,
		bson.M{
			"id": id,
		},
	).Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (b *Backend) UpdateCategory(c *models.Category) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CategoriesCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": c.Id,
		},
		bson.M{
			"$set": c,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) DeleteCategory(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(CategoriesCollection).DeleteOne(ctx,
		bson.M{
			"id": id,
		})
	if err != nil {
		return err
	}

	return nil
}

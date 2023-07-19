package mongo

import (
	"bar/internal/models"
	"time"

	"github.com/google/uuid"
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
			"id": uuid.MustParse(id),

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
		bson.M{
			"$set": c,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) MarkDeleteCategory(id, by string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CategoriesCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
		bson.M{
			"$set": bson.M{
				"deleted_at": time.Now().Unix(),
				"deleted_by": uuid.MustParse(by),
			},
		},
		options.FindOneAndUpdate().SetUpsert(false))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) UnMarkDeleteCategory(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CategoriesCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
		bson.M{
			"$set": bson.M{
				"deleted_at": nil,
				"deleted_by": nil,
			},
		},
		options.FindOneAndUpdate().SetUpsert(false))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) DeleteCategory(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CategoriesCollection).FindOneAndDelete(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) RestoreCategory(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(CategoriesCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
		bson.M{
			"$unset": bson.M{
				"deleted_at": "",
				"deleted_by": "",
			},
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

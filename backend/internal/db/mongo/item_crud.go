package mongo

import (
	"bar/internal/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) CreateItem(item *models.Item) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(ItemsCollection).InsertOne(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) GetItem(id string) (*models.Item, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var item models.Item
	err := b.db.Collection(ItemsCollection).FindOne(ctx,
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
	).Decode(&item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (b *Backend) UpdateItem(item *models.Item) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(ItemsCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": item.Id,

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
			"$set": item,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) MarkDeleteItem(id, by string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(ItemsCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) UnMarkDeleteItem(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(ItemsCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) DeleteItem(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(ItemsCollection).FindOneAndDelete(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) RestoreItem(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(ItemsCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) GetDeletedItems(page uint64, size uint64) ([]*models.Item, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var accs []*models.Item
	cursor, err := b.db.Collection(ItemsCollection).Find(ctx,
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

func (b *Backend) CountDeletedItems() (int64, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	count, err := b.db.Collection(ItemsCollection).CountDocuments(ctx, bson.M{
		"deleted_at": bson.M{
			"$ne": nil,
		},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

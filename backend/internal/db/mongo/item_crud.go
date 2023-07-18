package mongo

import (
	"bar/internal/models"
	"time"

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
			"id": id,
			"deleted_at": bson.M{
				"$exists": false,
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
			"deleted_at": bson.M{
				"$exists": false,
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

func (b *Backend) DeleteItem(id, by string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(ItemsCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": id,
		},
		bson.M{
			"$set": bson.M{
				"deleted_at": time.Now().Unix(),
				"deleted_by": by,
			},
		},
		options.FindOneAndUpdate().SetUpsert(false))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

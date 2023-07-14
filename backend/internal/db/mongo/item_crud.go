package mongo

import (
	"bar/internal/models"

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

func (b *Backend) DeleteItem(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(ItemsCollection).DeleteOne(ctx,
		bson.M{
			"id": id,
		})
	if err != nil {
		return err
	}

	return nil
}

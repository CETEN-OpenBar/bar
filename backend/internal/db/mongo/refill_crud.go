package mongo

import (
	"bar/internal/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) CreateRefill(refill *models.Refill) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	_, err := b.db.Collection(RefillsCollection).InsertOne(ctx, refill)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) GetRefill(id string) (*models.Refill, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var refill models.Refill
	err := b.db.Collection(RefillsCollection).FindOne(ctx,
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
	).Decode(&refill)
	if err != nil {
		return nil, err
	}

	return &refill, nil
}

func (b *Backend) UpdateRefill(refill *models.Refill) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(RefillsCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": refill.Id,

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
			"$set": refill,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) MarkDeleteRefill(id, by string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(RefillsCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) UnMarkDeleteRefill(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(RefillsCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) DeleteRefill(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(RefillsCollection).FindOneAndDelete(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) RestoreRefill(id string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	res := b.db.Collection(RefillsCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) GetDeletedRefills(page uint64, size uint64) ([]*models.Refill, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	var accs []*models.Refill
	cursor, err := b.db.Collection(RefillsCollection).Find(ctx,
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

func (b *Backend) CountDeletedRefills() (uint64, error) {
	ctx, cancel := b.GetContext()
	defer cancel()

	count, err := b.db.Collection(RefillsCollection).CountDocuments(ctx, bson.M{
		"deleted_at": bson.M{
			"$ne": nil,
		},
	})
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

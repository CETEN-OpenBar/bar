package mongo

import (
	"bar/internal/models"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

)

func (b *Backend) CreateStarring(ctx context.Context, starring *models.Starring) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()


	_, err := b.db.Collection(StarringCollection).InsertOne(ctx, starring)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) GetStarring(ctx context.Context, id string) (*models.Starring, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var starring models.Starring
	err := b.db.Collection(StarringCollection).FindOne(ctx,
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
	).Decode(&starring)
	if err != nil {
		return nil, err
	}

	return &starring, nil
}


func (b *Backend) UpdateStarring(ctx context.Context, starring *models.Starring) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res := b.db.Collection(StarringCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": starring.Id,

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
			"$set": starring,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) MarkDeleteStarring(ctx context.Context, id, by string) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res := b.db.Collection(StarringCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) UnMarkDeleteStarring(ctx context.Context, id string) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res := b.db.Collection(StarringCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
		bson.M{
			"$set": bson.M{
				"deleted_at":      nil,
				"deleted_by":      nil,
				"deleted_by_name": nil,
			},
		},
		options.FindOneAndUpdate().SetUpsert(false))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) DeleteStarring(ctx context.Context, id string) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res := b.db.Collection(StarringCollection).FindOneAndDelete(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) RestoreStarring(ctx context.Context, id string) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res := b.db.Collection(StarringCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) GetDeletedStarrings(ctx context.Context, page uint64, size uint64) ([]*models.Starring, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var accs []*models.Starring
	cursor, err := b.db.Collection(StarringCollection).Find(ctx,
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

func (b *Backend) CountDeletedStarrings(ctx context.Context) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	count, err := b.db.Collection(StarringCollection).CountDocuments(ctx, bson.M{
		"deleted_at": bson.M{
			"$ne": nil,
		},
	})
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

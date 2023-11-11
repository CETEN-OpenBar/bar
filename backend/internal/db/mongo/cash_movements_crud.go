package mongo

import (
	"bar/internal/models"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) CreateCashMovement(ctx context.Context, tx *models.CashMovement) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	tx.CashMovement.CreatedAt = uint64(time.Now().Unix())

	_, err := b.db.Collection(CashMovementsCollection).InsertOne(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) GetCashMovement(ctx context.Context, id string) (*models.CashMovement, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var tx models.CashMovement
	err := b.db.Collection(CashMovementsCollection).FindOne(ctx,
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
	).Decode(&tx)
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (b *Backend) UpdateCashMovement(ctx context.Context, tx *models.CashMovement) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res := b.db.Collection(CashMovementsCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": tx.Id,

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
			"$set": tx,
		},
		options.FindOneAndUpdate().SetUpsert(true))
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) MarkDeleteCashMovement(ctx context.Context, id, by string) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res := b.db.Collection(CashMovementsCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) UnMarkDeleteCashMovement(ctx context.Context, id string) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res := b.db.Collection(CashMovementsCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) DeleteCashMovement(ctx context.Context, id string) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res := b.db.Collection(CashMovementsCollection).FindOneAndDelete(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (b *Backend) RestoreCashMovement(ctx context.Context, id string) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res := b.db.Collection(CashMovementsCollection).FindOneAndUpdate(ctx,
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

func (b *Backend) GetDeletedCashMovements(ctx context.Context, page uint64, size uint64) ([]*models.CashMovement, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var accs []*models.CashMovement
	cursor, err := b.db.Collection(CashMovementsCollection).Find(ctx,
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

func (b *Backend) CountDeletedCashMovements(ctx context.Context) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	count, err := b.db.Collection(CashMovementsCollection).CountDocuments(ctx, bson.M{
		"deleted_at": bson.M{
			"$ne": nil,
		},
	})
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

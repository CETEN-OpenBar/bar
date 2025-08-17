package mongo

import (
	"bar/internal/models"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) CreateRemoteRefill(ctx context.Context, refill *models.RemoteRefill) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	refill.CreatedAt = uint64(time.Now().Unix())

	_, err := b.db.Collection(RemoteRefillsCollection).InsertOne(ctx, refill)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) GetRemoteRefill(ctx context.Context, id string) (*models.RemoteRefill, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var refill models.RemoteRefill
	err := b.db.Collection(RemoteRefillsCollection).FindOne(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
	).Decode(&refill)
	if err != nil {
		return nil, err
	}

	return &refill, nil
}

func (b *Backend) FindRemoteRefillForAccount(ctx context.Context, accountId string, checkoutIntentId int32) (*models.RemoteRefill, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var refill models.RemoteRefill
	err := b.db.Collection(RemoteRefillsCollection).FindOne(ctx,
		bson.M{
			"checkout_intent_id": checkoutIntentId,
			"account_id": uuid.MustParse(accountId),
		},
	).Decode(&refill)
	if err != nil {
		return nil, err
	}

	return &refill, nil
}

func (b *Backend) UpdateRemoteRefill(ctx context.Context, refill *models.RemoteRefill) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res := b.db.Collection(RemoteRefillsCollection).FindOneAndUpdate(ctx,
		bson.M{
			"id": refill.Id,
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

func (b *Backend) DeleteRemoteRefill(ctx context.Context, id string) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res := b.db.Collection(RemoteRefillsCollection).FindOneAndDelete(ctx,
		bson.M{
			"id": uuid.MustParse(id),
		},
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

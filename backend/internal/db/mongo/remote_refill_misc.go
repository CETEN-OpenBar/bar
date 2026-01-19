package mongo

import (
	"bar/autogen"
	"bar/internal/models"
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func (b *Backend) GetAllRemoteRefillsWithState(ctx context.Context, state autogen.RemoteRefillState) ([]*models.RemoteRefill, error) {
	
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	cursor, err := b.db.Collection(RemoteRefillsCollection).Find(ctx,
		bson.M{
			"state": state,
		},
		nil,
	)

	if err != nil {
		return nil, err
	}

	// Decode each refill
	var refills []*models.RemoteRefill
	if err := cursor.All(ctx, &refills); err != nil {
		return nil, err
	}

	return refills, nil
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

func (b *Backend) GetAllPendingRemoteRefillsForAccount(ctx context.Context, accountId string) ([]*models.RemoteRefill, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	cursor, err := b.db.Collection(RemoteRefillsCollection).Find(ctx,
		bson.M{
			"account_id": uuid.MustParse(accountId),
			"state": autogen.RemoteRefillStarted,
		},
	);

	if err != nil {
		return nil, err
	}

	var refills []*models.RemoteRefill
	if err := cursor.All(ctx, &refills); err != nil {
		return nil, err
	}

	return refills, nil
}

// Update a remote refill state atomically.
// The refill passed as a parameter is also updated if it was updated in the database
// Returns true if the refill was updated
func (b *Backend) UpdateRemoteRefillStateAtomic(ctx context.Context, refill *models.RemoteRefill, newState autogen.RemoteRefillState) (bool, error) {
	
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	res, err := b.db.Collection(RemoteRefillsCollection).UpdateOne(ctx,
		bson.M{
			"id": refill.Id,
			"state": refill.State,
		},
		bson.M{
			"$set": bson.M{
				"state": newState,
			},
		},
		nil,
	)

	if err != nil || res == nil {
		return false, err
	}

	if res.ModifiedCount != 1 {
		return false, nil
	}

	refill.State = newState
	return true, nil
}

func (b *Backend) GetAllRemoteRefills(ctx context.Context, page uint64, size uint64, accountName *string, state *autogen.RemoteRefillState, startAt, endAt uint64) ([]*models.RemoteRefill, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var refills []*models.RemoteRefill
	filter := bson.M{
		"created_at": bson.M{
			"$gte": startAt,
			"$lte": endAt,
		},
	}

	if accountName != nil {
		filter["account_name"] = bson.M{"$regex": *accountName, "$options": "i"}
	}

	if state != nil {
		filter["state"] = *state
	}

	cursor, err := b.db.Collection(RemoteRefillsCollection).Find(ctx, filter,
		options.Find().
			SetSkip(int64(page*size)).
			SetLimit(int64(size)).
			SetSort(bson.M{"created_at": -1}),
	)

	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &refills); err != nil {
		return nil, err
	}

	return refills, nil
}

func (b *Backend) CountAllRemoteRefills(ctx context.Context, accountName *string, state *autogen.RemoteRefillState, startAt, endAt uint64) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{
		"created_at": bson.M{
			"$gte": startAt,
			"$lte": endAt,
		},
	}

	if accountName != nil {
		filter["account_name"] = bson.M{"$regex": *accountName, "$options": "i"}
	}

	if state != nil {
		filter["state"] = *state
	}

	count, err := b.db.Collection(RemoteRefillsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}


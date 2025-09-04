package mongo

import (
	"bar/autogen"
	"bar/internal/models"
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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


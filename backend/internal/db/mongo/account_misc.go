package mongo

import (
	"bar/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (b *Backend) GetAccounts(ctx context.Context, role string, priceRole string, page uint64, size uint64, query string) ([]*models.Account, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	// Get "size" accounts from "page" using aggregation
	var accounts []*models.Account
	andConditions := []bson.M{
		{
			"$or": []bson.M{
				{"first_name": bson.M{"$regex": query, "$options": "i"}},
				{"last_name":  bson.M{"$regex": query, "$options": "i"}},
				{"email":      bson.M{"$regex": query, "$options": "i"}},
			},
		},
		{
			"$or": []bson.M{
				{"deleted_at": bson.M{"$exists": false}},
				{"deleted_at": nil},
			},
		},
	}

	if priceRole != "" {
		andConditions = append(andConditions, bson.M{
			"price_role": priceRole,
		})
	}

	if role != "" { 
		andConditions = append(andConditions, bson.M{
			"role": role,
		})
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"$and": andConditions,
			},
		},
		{
			"$skip": page * size,
		},
		{
			"$limit": size,
		},
	}

	cursor, err := b.db.Collection(AccountsCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	// Decode each account
	if err := cursor.All(ctx, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (b *Backend) CountAccounts(ctx context.Context, role string, priceRole string, query string) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()



	
	// Count all accounts
	andConditions := []bson.M{
		{
			"$or": []bson.M{
				{"first_name": bson.M{"$regex": query, "$options": "i"}},
				{"last_name":  bson.M{"$regex": query, "$options": "i"}},
				{"email":      bson.M{"$regex": query, "$options": "i"}},
			},
		},
		{
			"$or": []bson.M{
				{"deleted_at": bson.M{"$exists": false}},
				{"deleted_at": nil},
			},
		},
	}

	if priceRole != "" {
		andConditions = append(andConditions, bson.M{
			"price_role": priceRole,
		})
	}

	if role != "" {
		andConditions = append(andConditions, bson.M{
			"role": role, 
		})
	}

	count, err := b.db.Collection(AccountsCollection).CountDocuments(ctx, bson.M{
		"$and": andConditions,
	})


	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (b *Backend) GetAccountByCard(ctx context.Context, card string) (*models.Account, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	// Get account by card
	var account models.Account
	if err := b.db.Collection(AccountsCollection).FindOne(ctx, bson.M{"card_id": card}).Decode(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (b *Backend) GetAccountByGoogle(ctx context.Context, googleID string) (*models.Account, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	// Get account by card
	var account models.Account
	if err := b.db.Collection(AccountsCollection).FindOne(ctx, bson.M{"google_id": googleID}).Decode(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (b *Backend) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	// Get account by card
	var account models.Account
	if err := b.db.Collection(AccountsCollection).FindOne(ctx, bson.M{"email_address": email}).Decode(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (b *Backend) ListenForChanges(ctx context.Context, coll string, fn func(*mongo.ChangeStream)) error {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	// Watch for changes
	cs, err := b.db.Collection(coll).Watch(ctx, mongo.Pipeline{})
	if err != nil {
		return err
	}

	// Call function
	fn(cs)

	return nil
}

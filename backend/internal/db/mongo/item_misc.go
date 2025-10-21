package mongo

import (
	"bar/autogen"
	"bar/internal/models"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *Backend) GetItems(ctx context.Context, categoryID string, page, size uint64, state string, name string, fournisseur string, refBundle string, sort bool) ([]*models.Item, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var items []*models.Item

	matchFilter := bson.M{
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
	}
	if state != "" {
		matchFilter["state"] = state
		if state == string(autogen.ItemBuyable) {
						// Get seconds since day start
			t := time.Since(time.Now().Truncate(24 * time.Hour)).Seconds()
			// available_from <= t <= available_until or (available_from == nil && available_until == nil)
			matchFilter["$and"] = []bson.M{
				{
					"$or": []bson.M{
						{
							"available_from": bson.M{
								"$lte": t,
							},
						},
						{
							"available_from": nil,
						},
					},
				},
				{
					"$or": []bson.M{
						{
							"available_until": bson.M{
								"$gte": t,
							},
						},
						{
							"available_until": nil,
						},
					},
				},
			}
		}
	}
	if categoryID != "" {
		matchFilter["category_id"] = uuid.MustParse(categoryID)
	}
	if name != "" {
		matchFilter["name"] = bson.M{
			"$regex":   name,
			"$options": "i",
		}
	}
	if fournisseur != "" {
		matchFilter["fournisseur"] = fournisseur
	}
	if refBundle != "" {
		matchFilter["ref_bundle"] = refBundle
	}

	pipeline := []bson.M{
		{"$match": matchFilter},
	}

	if sort {
		pipeline = append(pipeline, bson.M{
			"$addFields": bson.M{
				"sortPriority": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{"$amount_left", 0}},
						2, 
						1,
					},
				},
			},
		})

		pipeline = append(pipeline, bson.M{
			"$sort": bson.D{
				{"sortPriority", 1},
				{"_id", 1}, 
			},
		})
	}
	pipeline = append(pipeline, bson.M{"$skip": int64(page * size)})
	pipeline = append(pipeline, bson.M{"$limit": int64(size)})

	if sort {
		pipeline = append(pipeline, bson.M{
			"$project": bson.M{"sortPriority": 0},
		})
	}

	cursor, err := b.db.Collection(ItemsCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (b *Backend) GetIncoherentItems(ctx context.Context, page, size uint64, categoryID string, state string, name string) ([]*models.Item, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	var items []*models.Item

	filter := bson.M{
		"$and": []bson.M{
			{
				"deleted_at": nil,
			},
			{
				"$or": []bson.M{
					{
						"$and": []bson.M{
							{
								"amount_left": bson.M{
									"$gt": 0,
								},
							},
							{
								"state": bson.M{
									"$ne": "buyable",
								},
							},
						},
					},
					{"prices.ceten": 0},
					{"prices.coutant": 0},
					{"prices.externe": 0},
					{"prices.menu": 0},
					{"prices.privilegies": 0},
					{"prices.staff_bar": 0},
				},
			},
		},
	}
	if state != "" {
		filter["state"] = state
		if state == string(autogen.ItemBuyable) {
			// Get seconds since day start
			t := time.Since(time.Now().Truncate(24 * time.Hour)).Seconds()
			// available_from <= t <= available_until or (available_from == nil && available_until == nil)
			filter["$and"] = []bson.M{
				{
					"$or": []bson.M{
						{
							"available_from": bson.M{
								"$lte": t,
							},
						},
						{
							"available_from": nil,
						},
					},
				},
				{
					"$or": []bson.M{
						{
							"available_until": bson.M{
								"$gte": t,
							},
						},
						{
							"available_until": nil,
						},
					},
				},
			}
		}
	}
	if categoryID != "" {
		filter["category_id"] = uuid.MustParse(categoryID)
	}
	if name != "" {
		filter["name"] = bson.M{
			"$regex":   name,
			"$options": "i",
		}
	}

	cursor, err := b.db.Collection(ItemsCollection).Find(ctx, filter, options.Find().SetSkip(int64(page*size)).SetLimit(int64(size)))
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (b *Backend) CountItems(ctx context.Context, categoryID string, state string, name string, fournisseur string, refBundle string) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{
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
	}
	if state != "" {
		filter["state"] = state
		if state == string(autogen.ItemBuyable) {
			t := time.Since(time.Now().Truncate(24 * time.Hour)).Seconds()
			filter["$and"] = []bson.M{
				{
					"$or": []bson.M{
						{
							"available_from": bson.M{
								"$lte": t,
							},
						},
						{
							"available_from": nil,
						},
					},
				},
				{
					"$or": []bson.M{
						{
							"available_until": bson.M{
								"$gte": t,
							},
						},
						{
							"available_until": nil,
						},
					},
				},
			}
		}
	}
	if categoryID != "" {
		filter["category_id"] = uuid.MustParse(categoryID)
	}
	if name != "" {
		filter["name"] = bson.M{
			"$regex":   name,
			"$options": "i",
		}
	}
	if fournisseur != "" {
		filter["fournisseur"] = fournisseur
	}
	if refBundle != "" {
		filter["ref_bundle"] = refBundle
	}

	count, err := b.db.Collection(ItemsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

func (b *Backend) CountIncoherentItems(ctx context.Context, categoryID string, state string, name string) (uint64, error) {
	ctx, cancel := b.TimeoutContext(ctx)
	defer cancel()

	filter := bson.M{
		"$and": []bson.M{
			{
				"deleted_at": nil,
			},
			{
				"$or": []bson.M{
					{
						"$and": []bson.M{
							{
								"amount_left": bson.M{
									"$gt": 0,
								},
							},
							{
								"state": bson.M{
									"$ne": "buyable",
								},
							},
						},
					},
					{"prices.ceten": 0},
					{"prices.coutant": 0},
					{"prices.externe": 0},
					{"prices.menu": 0},
					{"prices.privilegies": 0},
					{"prices.staff_bar": 0},
				},
			},
		},
	}

	if state != "" {
		filter["state"] = state
		if state == string(autogen.ItemBuyable) {
			t := time.Since(time.Now().Truncate(24 * time.Hour)).Seconds()
			filter["$and"] = []bson.M{
				{
					"$or": []bson.M{
						{
							"available_from": bson.M{
								"$lte": t,
							},
						},
						{
							"available_from": nil,
						},
					},
				},
				{
					"$or": []bson.M{
						{
							"available_until": bson.M{
								"$gte": t,
							},
						},
						{
							"available_until": nil,
						},
					},
				},
			}
		}
	}
	if categoryID != "" {
		filter["category_id"] = uuid.MustParse(categoryID)
	}
	if name != "" {
		filter["name"] = bson.M{
			"$regex":   name,
			"$options": "i",
		}
	}

	count, err := b.db.Collection(ItemsCollection).CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}

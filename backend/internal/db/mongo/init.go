package mongo

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Collections = []string{
		"transactions",
		"refills",
		"remote_refills",
		"starrings",
		"items",
		"categories",
		"carousel_texts",
		"carousel_images",
		"accounts",
		"restocks",
		"cash_movements",
	}

	// Add indexes to index "id" which is text & unique
	indexes = map[string][]mongo.IndexModel{
		"transactions": {
			// index descending by created_at
			mongo.IndexModel{
				Keys: bson.M{
					"created_at": -1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"id": bson.M{
						"$exists": true,
					},
				}),
			},
		},
		"refills": {
			mongo.IndexModel{
				Keys: bson.M{
					"created_at": 1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"id": bson.M{
						"$exists": true,
					},
				}),
			}},
		"starrings": {
			mongo.IndexModel{
				Keys: bson.M{
					"created_at": 1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"id": bson.M{
						"$exists": true,
					},
				}),
			}},
		"remote_refills": {
			mongo.IndexModel{
				Keys: bson.M{
					"created_at": 1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"checkout_intent_id": 1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"order_id": 1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"id": bson.M{
						"$exists": true,
					},
				}),
			}},
		"items": {
			mongo.IndexModel{
				Keys: bson.M{
					"created_at": 1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"id": bson.M{
						"$exists": true,
					},
				}),
			}},
		"categories": {
			mongo.IndexModel{
				Keys: bson.M{
					"created_at": 1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"id": bson.M{
						"$exists": true,
					},
				}),
			}},
		"carousel_texts": {
			mongo.IndexModel{
				Keys: bson.M{
					"created_at": 1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"id": bson.M{
						"$exists": true,
					},
				}),
			}},
		"carousel_images": {
			mongo.IndexModel{
				Keys: bson.M{
					"created_at": 1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"id": bson.M{
						"$exists": true,
					},
				}),
			}},
		"accounts": {
			mongo.IndexModel{
				Keys: bson.M{
					"created_at": 1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"id": bson.M{
						"$exists": true,
					},
				}),
			},
			mongo.IndexModel{
				// Avoid duplicates except empty values
				Keys: bson.M{
					"card_id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"$and": []bson.M{
						{
							"card_id": bson.M{
								"$exists": true,
							},
						},
						{
							"card_id": bson.M{
								"$type": "string",
							},
						},
					},
				}),
			},
			mongo.IndexModel{
				Keys: bson.M{
					"google_id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"$and": []bson.M{
						{
							"google_id": bson.M{
								"$exists": true,
							},
						},
						{
							"google_id": bson.M{
								"$type": "string",
							},
						},
					},
				}),
			},
			mongo.IndexModel{
				Keys: bson.M{
					"email_address": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"$and": []bson.M{
						{
							"email_address": bson.M{
								"$exists": true,
							},
						},
						{
							"email_address": bson.M{
								"$type": "string",
							},
						},
					},
				}),
			},
		},
		"restocks": {
			mongo.IndexModel{
				Keys: bson.M{
					"created_at": 1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"id": bson.M{
						"$exists": true,
					},
				}),
			},
		},
		"cash_movements": {
			mongo.IndexModel{
				Keys: bson.M{
					"created_at": 1,
				},
			},
			mongo.IndexModel{
				Keys: bson.M{
					"id": 1,
				},
				Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
					"id": bson.M{
						"$exists": true,
					},
				}),
			},
		},
	}

	TransactionsCollection   = "transactions"
	RefillsCollection        = "refills"
	RemoteRefillsCollection  = "remote_refills"
	StarringCollection       = "starrings"
	ItemsCollection          = "items"
	CategoriesCollection     = "categories"
	CarouselTextsCollection  = "carousel_texts"
	CarouselImagesCollection = "carousel_images"
	AccountsCollection       = "accounts"
	RestocksCollection       = "restocks"
	CashMovementsCollection  = "cash_movements"
)

func (b *Backend) CreateCollections() error {
	ctx, cancel := b.GetContext()
	defer cancel()

	for _, collection := range Collections {
		b.db.CreateCollection(ctx, collection)

		if err := b.CreateIndexes(collection); err != nil {
			logrus.Error(err)
			continue
		}
	}

	return nil
}

func (b *Backend) CreateIndexes(collection string) error {
	ctx, cancel := b.GetContext()
	defer cancel()

	v, ok := indexes[collection]
	if !ok {
		return nil
	}

	_, err := b.db.Collection(collection).Indexes().CreateMany(ctx, v)
	if err != nil {
		return err
	}

	return nil
}

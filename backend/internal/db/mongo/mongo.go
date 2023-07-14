package mongo

import (
	"bar/internal/db"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Backend struct {
	db   *mongo.Database
	opts *db.DatabaseOptions
}

func NewMongoBackend(opts *db.DatabaseOptions) db.DBackend {
	return &Backend{
		opts: opts,
	}
}

func (b *Backend) Connect() error {
	ctx, cancel := b.GetContext()
	defer cancel()

	// Connect to mongodb
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(b.opts.ConnString))
	if err != nil {
		return err
	}

	b.db = client.Database(b.opts.Database)

	// Starts all inits

	return nil
}

func (b *Backend) Disconnect() error {
	ctx, cancel := b.GetContext()
	defer cancel()

	return b.db.Client().Disconnect(ctx)
}

func (b *Backend) GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), b.opts.QueryTimeout)
}

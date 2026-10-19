// Command seed fills the local development database with sample data.
// Run it with `task seed` from the repo root, or `go run ./cmd/seed` here.
package main

import (
	"bar/internal/config"
	"bar/internal/db"
	"bar/internal/db/mongo"
	"bar/internal/seed"
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	c := config.GetConfig()

	if !seed.IsLocalURI(c.MongoConfig.ConnectionURI) {
		fmt.Fprintf(os.Stderr, "refusing to seed non-local database: %s\n", c.MongoConfig.ConnectionURI)
		os.Exit(1)
	}

	opts := db.NewDatabaseOptions(c.MongoConfig.ConnectionURI, c.MongoConfig.Database, time.Millisecond*time.Duration(c.MongoConfig.Timeout))
	b := mongo.NewMongoBackend(opts)
	if err := b.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to %s: %v\n", c.MongoConfig.ConnectionURI, err)
		os.Exit(1)
	}
	defer b.Disconnect()

	n, err := seed.Run(context.Background(), b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed failed after %d documents: %v\n", n, err)
		os.Exit(1)
	}
	fmt.Printf("seeded %d documents\n", n)
}

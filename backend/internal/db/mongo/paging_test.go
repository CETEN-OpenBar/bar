package mongo

import (
	"bar/internal/db"
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// A size of 0 must not send $limit: 0, which Mongo rejects
func TestAggregatedPagesWithoutLimit(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	queries := map[string]func(b *Backend, size uint64) error{
		"accounts": func(b *Backend, size uint64) error {
			_, err := b.GetAccounts(context.Background(), "", "", 1, size, "")
			return err
		},
		"items": func(b *Backend, size uint64) error {
			_, err := b.GetItems(context.Background(), "", 1, size, "", "", "", "", true)
			return err
		},
	}

	for name, query := range queries {
		for _, size := range []uint64{0, 5} {
			mt.Run(fmt.Sprintf("%s/size=%d", name, size), func(mt *mtest.T) {
				backend := &Backend{db: mt.DB, opts: &db.DatabaseOptions{QueryTimeout: time.Second}}
				mt.AddMockResponses(mtest.CreateCursorResponse(0, mt.DB.Name()+".x", mtest.FirstBatch))
				if err := query(backend, size); err != nil {
					mt.Fatal(err)
				}

				stages, err := mt.GetStartedEvent().Command.Lookup("pipeline").Array().Values()
				if err != nil {
					mt.Fatal(err)
				}
				var skip, limit bool
				for _, stage := range stages {
					doc := stage.Document()
					if _, err := doc.LookupErr("$skip"); err == nil {
						skip = true
					}
					if _, err := doc.LookupErr("$limit"); err == nil {
						limit = true
					}
				}
				if want := size > 0; skip != want || limit != want {
					mt.Errorf("$skip sent: %v, $limit sent: %v, want %v for both", skip, limit, want)
				}
			})
		}
	}
}

package mongo

import (
	"bar/autogen"
	"bar/internal/db"
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestAllTransactionsQueriesPreserveSearchScope(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	for _, state := range []autogen.TransactionState{
		autogen.TransactionStarted,
		autogen.TransactionTakenCareOf,
		autogen.TransactionFinished,
		autogen.TransactionCanceled,
	} {
		for _, hideRemotes := range []bool{false, true} {
			mt.Run(fmt.Sprintf("%s/hide_remotes=%t", state, hideRemotes), func(mt *mtest.T) {
				backend := &Backend{db: mt.DB, opts: &db.DatabaseOptions{QueryTimeout: time.Second}}
				namespace := mt.DB.Name() + "." + TransactionsCollection
				mt.AddMockResponses(
					mtest.CreateCursorResponse(0, namespace, mtest.FirstBatch),
					mtest.CreateCursorResponse(0, namespace, mtest.FirstBatch, bson.D{{Key: "n", Value: int32(75)}}),
				)

				// Combine the name, remote, item and date filters on a later page.
				const itemID = "704cc10b-9306-4627-9fef-e94c2d8fa4b2"
				_, err := backend.GetAllTransactions(context.Background(), 2, 30, string(state), true, "Alice (bar)+", hideRemotes, 100, 200, itemID)
				if err != nil {
					mt.Fatal(err)
				}
				count, err := backend.CountAllTransactions(context.Background(), string(state), true, "Alice (bar)+", hideRemotes, 100, 200, itemID)
				if err != nil {
					mt.Fatal(err)
				}
				if count != 75 {
					mt.Fatalf("count = %d, want 75", count)
				}

				find := mt.GetStartedEvent()
				countCommand := mt.GetStartedEvent()
				if find == nil || countCommand == nil || find.CommandName != "find" || countCommand.CommandName != "aggregate" {
					mt.Fatal("expected a find followed by the count aggregation")
				}
				filter := find.Command.Lookup("filter").Document()
				pipeline := countCommand.Command.Lookup("pipeline").Array()
				match := pipeline.Index(0).Value().Document().Lookup("$match").Document()
				assertSameBSONDocument(mt.T, filter, match)
				if got := filter.Lookup("state").StringValue(); got != string(state) {
					mt.Fatalf("search state = %q, want %q", got, state)
				}

				nameConditions, err := filter.Lookup("$or").Array().Values()
				if err != nil || len(nameConditions) != 2 {
					mt.Fatal("name and nickname alternatives must survive the remote filter")
				}
				for i, field := range []string{"account_name", "account_nick_name"} {
					condition := nameConditions[i].Document().Lookup(field).Document()
					pattern := condition.Lookup("$regex").StringValue()
					flags := condition.Lookup("$options").StringValue()
					re, err := regexp.Compile("(?" + flags + ")" + pattern)
					if err != nil {
						mt.Fatal(err)
					}
					if !re.MatchString("Mme ALICE (BAR)+ Dupont") || re.MatchString("Alice barrr") {
						mt.Fatalf("%s must match a literal, case-insensitive name substring", field)
					}
				}
				remoteCondition, err := filter.LookupErr("is_remote")
				if hideRemotes {
					// $ne true accepts false, null and missing legacy values in MongoDB.
					if err != nil || !remoteCondition.Document().Lookup("$ne").Boolean() {
						mt.Fatal("remote transactions must be excluded without excluding legacy transactions")
					}
				} else if err == nil {
					mt.Fatal("remote transactions must remain eligible when the toggle is off")
				}
				if got := find.Command.Lookup("skip").AsInt64(); got != 60 {
					mt.Fatalf("skip = %d, want 60", got)
				}
				if got := find.Command.Lookup("limit").AsInt64(); got != 30 {
					mt.Fatalf("limit = %d, want 30", got)
				}
				if got := find.Command.Lookup("sort").Document().Lookup("created_at").AsInt64(); got != -1 {
					mt.Fatalf("created_at sort = %d, want -1", got)
				}
				stages, err := pipeline.Values()
				if err != nil {
					mt.Fatal(err)
				}
				for _, stage := range stages {
					for _, field := range []string{"$skip", "$limit"} {
						if _, err := stage.Document().LookupErr(field); err == nil {
							mt.Fatalf("count must cover all matching transactions, found %s", field)
						}
					}
				}
			})
		}
	}
}

func TestAllTransactionsQueryWithoutOptionalFilters(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("all transactions", func(mt *mtest.T) {
		backend := &Backend{db: mt.DB, opts: &db.DatabaseOptions{QueryTimeout: time.Second}}
		mt.AddMockResponses(mtest.CreateCursorResponse(0, mt.DB.Name()+"."+TransactionsCollection, mtest.FirstBatch))
		if _, err := backend.GetAllTransactions(context.Background(), 0, 30, "", false, "", false, 0, 0, ""); err != nil {
			mt.Fatal(err)
		}
		command := mt.GetStartedEvent().Command
		filter, err := command.Lookup("filter").Document().Elements()
		if err != nil || len(filter) != 0 {
			mt.Fatal("an unfiltered request must remain supported")
		}
	})
}

func TestTransactionsCreateIndexForStateAndNewestFirst(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("compound index", func(mt *mtest.T) {
		backend := &Backend{db: mt.DB, opts: &db.DatabaseOptions{QueryTimeout: time.Second}}
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		if err := backend.CreateIndexes(TransactionsCollection); err != nil {
			mt.Fatal(err)
		}
		createdIndexes, err := mt.GetStartedEvent().Command.Lookup("indexes").Array().Values()
		if err != nil {
			mt.Fatal(err)
		}
		for _, index := range createdIndexes {
			keys, err := index.Document().Lookup("key").Document().Elements()
			if err == nil && len(keys) == 2 && keys[0].Key() == "state" && keys[0].Value().AsInt64() == 1 && keys[1].Key() == "created_at" && keys[1].Value().AsInt64() == -1 {
				return
			}
		}
		mt.Fatal("missing ordered compound index on state then newest created_at")
	})
}

func assertSameBSONDocument(t *testing.T, first, second bson.Raw) {
	t.Helper()
	var firstMap, secondMap bson.M
	if err := bson.Unmarshal(first, &firstMap); err != nil {
		t.Fatal(err)
	}
	if err := bson.Unmarshal(second, &secondMap); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(firstMap, secondMap) {
		t.Fatalf("find and count filters differ: %v / %v", firstMap, secondMap)
	}
}

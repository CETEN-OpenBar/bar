//go:build integration

package mongo

import (
	"bar/autogen"
	"bar/internal/db"
	"bar/internal/models"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Set by TestMain, shared by every test in the package
var connString string

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Same image as docker-compose.yml. Change streams and transactions only
	// work on a replica set.
	container, err := mongodb.Run(ctx, "mongo:7.0", mongodb.WithReplicaSet("rs0"))
	if err != nil {
		log.Fatalf("failed to start MongoDB: %v", err)
	}
	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get the MongoDB host: %v", err)
	}
	// Docker often publishes ports on ::1 without forwarding them
	if host == "localhost" {
		host = "127.0.0.1"
	}
	port, err := container.MappedPort(ctx, "27017/tcp")
	if err != nil {
		log.Fatalf("failed to get the MongoDB port: %v", err)
	}
	// The replica set advertises the container IP, which the host can't
	// always reach (Docker Desktop, some firewalls). Talk to the published
	// port directly, like mongo-express does in docker-compose.yml.
	connString = fmt.Sprintf("mongodb://%s/?directConnection=true", net.JoinHostPort(host, port.Port()))

	code := m.Run()

	if err := testcontainers.TerminateContainer(container); err != nil {
		log.Printf("failed to stop MongoDB: %v", err)
	}
	os.Exit(code)
}

// newBackend connects to a database of its own, dropped when the test ends
func newBackend(t *testing.T) *Backend {
	t.Helper()

	name := strings.NewReplacer("/", "_", " ", "_", ".", "_").Replace(t.Name())
	b := NewMongoBackend(db.NewDatabaseOptions(connString, name, 10*time.Second)).(*Backend)
	if err := b.Connect(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := b.db.Drop(context.Background()); err != nil {
			t.Error(err)
		}
		if err := b.Disconnect(); err != nil {
			t.Error(err)
		}
	})
	return b
}

func createAccounts(t *testing.T, b *Backend, n int) []*models.Account {
	t.Helper()

	accounts := make([]*models.Account, n)
	for i := range accounts {
		accounts[i] = &models.Account{Account: autogen.Account{
			Id:           uuid.New(),
			FirstName:    fmt.Sprintf("User%d", i),
			LastName:     "Test",
			EmailAddress: fmt.Sprintf("user%d@example.com", i),
			Role:         autogen.AccountStudent,
			PriceRole:    autogen.AccountPriceCeten,
			State:        autogen.AccountOK,
		}}
		if err := b.CreateAccount(context.Background(), accounts[i]); err != nil {
			t.Fatal(err)
		}
	}
	return accounts
}

func accountNames(accounts []*models.Account) []string {
	names := make([]string, len(accounts))
	for i, acc := range accounts {
		names[i] = acc.FirstName
	}
	return names
}

func TestGetAccountsPages(t *testing.T) {
	b := newBackend(t)
	ctx := context.Background()

	accounts := createAccounts(t, b, 5)
	if err := b.MarkDeleteAccount(ctx, accounts[4].Id.String(), uuid.NewString()); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		page, size uint64
		want       []string
	}{
		{"no limit returns everything", 0, 0, []string{"User0", "User1", "User2", "User3"}},
		{"first page", 0, 3, []string{"User0", "User1", "User2"}},
		{"last page", 1, 3, []string{"User3"}},
		{"past the end", 2, 3, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := b.GetAccounts(ctx, "", "", tt.page, tt.size, "")
			if err != nil {
				t.Fatal(err)
			}
			if names := accountNames(got); strings.Join(names, ",") != strings.Join(tt.want, ",") {
				t.Errorf("GetAccounts() = %v, want %v", names, tt.want)
			}
		})
	}

	count, err := b.CountAccounts(ctx, "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Errorf("CountAccounts() = %d, want 4", count)
	}
}

func TestGetItemsPutsSoldOutLast(t *testing.T) {
	b := newBackend(t)
	ctx := context.Background()

	category := uuid.New()
	for i, left := range []uint64{0, 3, 0, 7} {
		item := &models.Item{Item: autogen.Item{
			Id:         uuid.New(),
			CategoryId: category,
			Name:       fmt.Sprintf("Item%d", i),
			AmountLeft: left,
			State:      autogen.ItemBuyable,
		}}
		if err := b.CreateItem(ctx, item); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name       string
		page, size uint64
		want       []string
	}{
		{"no limit returns everything", 0, 0, []string{"Item1", "Item3", "Item0", "Item2"}},
		{"second page", 1, 2, []string{"Item0", "Item2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := b.GetItems(ctx, category.String(), tt.page, tt.size, "", "", "", "", true)
			if err != nil {
				t.Fatal(err)
			}
			names := make([]string, len(items))
			for i, item := range items {
				names[i] = item.Name
			}
			if strings.Join(names, ",") != strings.Join(tt.want, ",") {
				t.Errorf("GetItems() = %v, want %v", names, tt.want)
			}
		})
	}
}

// WithTransaction backs the money paths: nothing may be kept when they fail
func TestWithTransactionRollsBack(t *testing.T) {
	b := newBackend(t)
	ctx := context.Background()

	acc := createAccounts(t, b, 1)[0]
	failure := errors.New("payment refused")

	_, err := b.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		acc.Balance = 1000
		if err := b.UpdateAccount(ctx, acc); err != nil {
			return nil, err
		}
		return nil, failure
	})
	if !errors.Is(err, failure) {
		t.Fatalf("WithTransaction() error = %v, want %v", err, failure)
	}

	got, err := b.GetAccount(ctx, acc.Id.String())
	if err != nil {
		t.Fatal(err)
	}
	if got.Balance != 0 {
		t.Errorf("balance = %d after a failed transaction, want 0", got.Balance)
	}
}

// The account watcher relies on this to push balance changes to the kiosks
func TestListenForChangesSeesUpdates(t *testing.T) {
	b := newBackend(t)
	ctx := context.Background()

	acc := createAccounts(t, b, 1)[0]

	watching := make(chan struct{})
	changed := make(chan bson.M, 1)
	errs := make(chan error, 1)
	go func() {
		errs <- b.ListenForChanges(ctx, AccountsCollection, func(cs *mongo.ChangeStream) {
			close(watching)
			if !cs.Next(ctx) {
				return
			}
			var event bson.M
			if err := cs.Decode(&event); err == nil {
				changed <- event
			}
		})
	}()

	select {
	case <-watching:
	case err := <-errs:
		t.Fatalf("ListenForChanges() = %v before watching", err)
	case <-time.After(10 * time.Second):
		t.Fatal("change stream never started")
	}

	acc.Balance = 500
	if err := b.UpdateAccount(ctx, acc); err != nil {
		t.Fatal(err)
	}

	select {
	case event := <-changed:
		if op := event["operationType"]; op != "update" && op != "replace" {
			t.Errorf("operationType = %v, want an update", op)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("no change event for the account update")
	}
	if err := <-errs; err != nil {
		t.Errorf("ListenForChanges() = %v", err)
	}
}

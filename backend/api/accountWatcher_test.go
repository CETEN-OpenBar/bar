package api

import (
	"bar/internal/db"
	"bar/internal/models"
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// watchBackend feeds AccountWatcher a mocked change stream once start is closed
type watchBackend struct {
	db.DBackend
	mt     *mtest.T
	start  chan struct{}
	events []primitive.ObjectID
	done   chan struct{}
}

func newWatchBackend(mt *mtest.T, events ...primitive.ObjectID) *watchBackend {
	b := &watchBackend{mt: mt, start: make(chan struct{}), events: events, done: make(chan struct{})}
	// The watcher wakes every listener on its way out: let it finish before
	// the next test registers its own
	mt.T.Cleanup(func() {
		waitClosed(mt.T, b.done, "stream never stopped")
		waitWatcherStopped(mt.T)
	})
	return b
}

func (b *watchBackend) ListenForChanges(ctx context.Context, coll string, fn func(*mongo.ChangeStream)) error {
	defer close(b.done)
	<-b.start

	docs := make([]bson.D, len(b.events))
	for i, id := range b.events {
		docs[i] = bson.D{
			{Key: "_id", Value: bson.D{{Key: "_data", Value: strconv.Itoa(i)}}},
			{Key: "operationType", Value: "update"},
			{Key: "documentKey", Value: bson.D{{Key: "_id", Value: id}}},
		}
	}
	b.mt.AddMockResponses(mtest.CreateCursorResponse(0, b.mt.DB.Name()+"."+coll, mtest.FirstBatch, docs...))
	cs, err := b.mt.Coll.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		return err
	}
	defer cs.Close(context.Background())
	fn(cs)
	return nil
}

func (b *watchBackend) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	return &models.Account{}, nil
}

// watch runs WatchAccount for account and closes the returned channel when it returns
func watch(s *Server, account *models.Account, ctx context.Context) <-chan struct{} {
	req := httptest.NewRequest(http.MethodGet, "/account/watch", nil).WithContext(ctx)
	c := echo.New().NewContext(req, httptest.NewRecorder())
	c.Set("userLogged", true)
	c.Set("userAccount", account)

	done := make(chan struct{})
	go func() {
		defer close(done)
		s.WatchAccount(c)
	}()
	return done
}

// Upper bound for things that should happen right away, generous for slow CI runners
const timeout = 5 * time.Second

func listenerCount(id primitive.ObjectID) int {
	mutex.Lock()
	defer mutex.Unlock()
	return len(listeners[id])
}

func waitForListeners(t *testing.T, id primitive.ObjectID, n int) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for listenerCount(id) != n {
		if time.Now().After(deadline) {
			t.Fatalf("got %d listeners, want %d", listenerCount(id), n)
		}
		time.Sleep(time.Millisecond)
	}
}

func waitWatcherStopped(t *testing.T) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		mutex.Lock()
		running := watcherRunning
		mutex.Unlock()
		if !running {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("watcher still running")
		}
		time.Sleep(time.Millisecond)
	}
}

func waitClosed(t *testing.T, ch <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(timeout):
		t.Fatal(what)
	}
}

func TestWatchAccountReturnsWhenClientLeaves(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("leave", func(mt *mtest.T) {
		b := newWatchBackend(mt)
		defer func() { close(b.start); <-b.done }()

		account := &models.Account{PrivateID: primitive.NewObjectID()}
		ctx, cancel := context.WithCancel(context.Background())
		done := watch(&Server{DBackend: b}, account, ctx)
		waitForListeners(mt.T, account.PrivateID, 1)

		cancel()
		waitClosed(mt.T, done, "WatchAccount still blocked after the client left")
		if n := listenerCount(account.PrivateID); n != 0 {
			mt.Errorf("%d listeners left behind, want 0", n)
		}
	})
}

// Two changes to the same account in a row: the second one comes in before
// the woken handler has unregistered its listener.
func TestAccountWatcherSurvivesBackToBackChanges(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("back to back", func(mt *mtest.T) {
		account := &models.Account{PrivateID: primitive.NewObjectID()}
		b := newWatchBackend(mt, account.PrivateID, account.PrivateID)

		done := watch(&Server{DBackend: b}, account, context.Background())
		waitForListeners(mt.T, account.PrivateID, 1)

		// Hold the handler between waking up and unregistering
		mutex.Lock()
		close(b.start)
		time.Sleep(50 * time.Millisecond)
		mutex.Unlock()

		waitClosed(mt.T, done, "WatchAccount was never notified")
		waitClosed(mt.T, b.done, "watcher stuck notifying a listener nobody reads")
	})
}

// Run with -race: the watcher walks listeners while handlers register
func TestAccountWatcherWhileHandlersRegister(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("register", func(mt *mtest.T) {
		changed := primitive.NewObjectID()
		events := make([]primitive.ObjectID, 200)
		for i := range events {
			events[i] = changed
		}
		b := newWatchBackend(mt, events...)
		s := &Server{DBackend: b}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// Start the watcher first so later handlers don't race to start it
		first := &models.Account{PrivateID: primitive.NewObjectID()}
		watch(s, first, ctx)
		waitForListeners(mt.T, first.PrivateID, 1)

		var wg sync.WaitGroup
		close(b.start)
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				watch(s, &models.Account{PrivateID: primitive.NewObjectID()}, ctx)
			}()
		}
		wg.Wait()
		waitClosed(mt.T, b.done, "watcher did not finish the stream")
	})
}

func TestWatchAccountWokenWhenStreamStops(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("stops", func(mt *mtest.T) {
		// No events: the stream ends as soon as it starts
		b := newWatchBackend(mt)
		account := &models.Account{PrivateID: primitive.NewObjectID()}

		done := watch(&Server{DBackend: b}, account, context.Background())
		waitForListeners(mt.T, account.PrivateID, 1)
		close(b.start)

		waitClosed(mt.T, done, "WatchAccount still waiting on a stream that stopped")
	})
}

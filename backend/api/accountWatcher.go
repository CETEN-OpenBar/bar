package api

import (
	"bar/autogen"
	"context"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// mutex guards listeners and watcherRunning
var mutex sync.Mutex
var listeners = make(map[primitive.ObjectID][]chan int)
var watcherRunning = false

// notify wakes the listeners of an account. Listeners are buffered and only
// read once, so one that already has a pending notification is skipped
// instead of blocking the watcher.
func notify(ls []chan int) {
	for _, listener := range ls {
		select {
		case listener <- 1:
		default:
		}
	}
}

func (s *Server) AccountWatcher() {
	// If already running, return
	mutex.Lock()
	if watcherRunning {
		mutex.Unlock()
		return
	}
	watcherRunning = true
	mutex.Unlock()

	err := s.DBackend.ListenForChanges(context.TODO(), "accounts", func(changeStream *mongo.ChangeStream) {
		for changeStream.Next(context.TODO()) {
			// Get account id from change stream
			var accountChange struct {
				DocumentKey struct {
					Id primitive.ObjectID `bson:"_id"`
				} `bson:"documentKey"`
			}
			err := changeStream.Decode(&accountChange)
			if err != nil {
				logrus.Error(err)
				continue
			}

			// Notify listeners
			mutex.Lock()
			notify(listeners[accountChange.DocumentKey.Id])
			mutex.Unlock()
		}
	})
	if err != nil {
		logrus.Error(err)
	}

	// The stream is gone: wake everyone so clients poll again, which restarts it
	mutex.Lock()
	watcherRunning = false
	for _, ls := range listeners {
		notify(ls)
	}
	mutex.Unlock()
}

// (GET /account/watch)
func (s *Server) WatchAccount(c echo.Context) error {
	account, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	// Add a listener to the account
	listener := make(chan int, 1)
	mutex.Lock()
	listeners[account.PrivateID] = append(listeners[account.PrivateID], listener)
	mutex.Unlock()

	// Make sure the stream is up, after registering so a stream that stops
	// right away still wakes this listener
	go s.AccountWatcher()

	// Wait for a notification, or for the client to leave
	ctx := c.Request().Context()
	select {
	case <-listener:
	case <-ctx.Done():
	}

	mutex.Lock()
	// Remove the listener
	for i, l := range listeners[account.PrivateID] {
		if l == listener {
			listeners[account.PrivateID] = append(listeners[account.PrivateID][:i], listeners[account.PrivateID][i+1:]...)
			break
		}
	}
	mutex.Unlock()

	// Nobody is left to answer
	if ctx.Err() != nil {
		return nil
	}

	// Return account
	account, err = s.DBackend.GetAccount(ctx, account.Id.String())
	if err != nil {
		return err
	}

	r := autogen.WatchAccount200JSONResponse{
		Account: &account.Account,
	}
	return r.VisitWatchAccountResponse(c.Response())
}

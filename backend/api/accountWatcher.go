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

var mutex sync.Mutex
var listeners = make(map[primitive.ObjectID][]chan int)
var watcherRunning = false

func (s *Server) AccountWatcher() {
	// If already running, return
	if watcherRunning {
		return
	}
	watcherRunning = true

	s.DBackend.ListenForChanges(context.TODO(), "accounts", func(changeStream *mongo.ChangeStream) {
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
			for _, listener := range listeners[accountChange.DocumentKey.Id] {
				listener <- 1
			}
		}
	})

	watcherRunning = false
}

// (GET /account/watch)
func (s *Server) WatchAccount(c echo.Context) error {
	account, err := MustGetUser(c)
	if err != nil {
		return nil
	}

	// Make sure the stream is up
	go s.AccountWatcher()

	// Add a listener to the account
	listener := make(chan int)
	mutex.Lock()
	listeners[account.PrivateID] = append(listeners[account.PrivateID], listener)
	mutex.Unlock()

	// Wait for a notification
	<-listener

	mutex.Lock()
	// Remove the listener
	for i, l := range listeners[account.PrivateID] {
		if l == listener {
			listeners[account.PrivateID] = append(listeners[account.PrivateID][:i], listeners[account.PrivateID][i+1:]...)
			break
		}
	}
	mutex.Unlock()

	// Return account
	account, err = s.DBackend.GetAccount(c.Request().Context(), account.Id.String())
	if err != nil {
		return err
	}

	r := autogen.WatchAccount200JSONResponse{
		Account: &account.Account,
	}
	return r.VisitWatchAccountResponse(c.Response())
}

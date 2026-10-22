package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/yyewolf/acr122u"
)

type fakeCard struct{}

func (fakeCard) Reader() string                  { return "fake" }
func (fakeCard) Status() (acr122u.Status, error) { return acr122u.Status{}, nil }
func (fakeCard) UID() []byte                     { return []byte{0xde, 0xad, 0xbe, 0xef} }

// Kiosks connect and disconnect while cards are tapped
func TestServeCardWhileClientsReconnect(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(cardReader))
	defer srv.Close()
	url := "ws" + strings.TrimPrefix(srv.URL, "http")

	logrus.SetLevel(logrus.WarnLevel)
	h := &handler{logrus.StandardLogger()}

	stop := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				c, _, err := websocket.DefaultDialer.Dial(url, nil)
				if err != nil {
					t.Error(err)
					return
				}
				c.Close()
			}
		}()
	}

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		h.ServeCard(fakeCard{})
	}
	close(stop)
	wg.Wait()
}

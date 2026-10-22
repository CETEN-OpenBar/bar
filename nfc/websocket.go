package main

import (
	"flag"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var addr = flag.String("addr", "localhost:3737", "http service address")

// listenersMu guards listeners: the card reader and every connection's
// handler run on their own goroutines
var listenersMu sync.Mutex
var listeners = make(map[*websocket.Conn]bool)

// connectedListeners returns a snapshot so writes happen without the lock
func connectedListeners() []*websocket.Conn {
	listenersMu.Lock()
	defer listenersMu.Unlock()
	conns := make([]*websocket.Conn, 0, len(listeners))
	for l := range listeners {
		conns = append(conns, l)
	}
	return conns
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func cardReader(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	listenersMu.Lock()
	listeners[c] = true
	listenersMu.Unlock()
	defer c.Close()
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
	listenersMu.Lock()
	delete(listeners, c)
	listenersMu.Unlock()
}

func server() {
	http.HandleFunc("/", cardReader)
	http.HandleFunc("/config.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "config.json")
	})
	logrus.Info("listening on ", *addr)
	logrus.Fatal(http.ListenAndServe(*addr, nil))
}

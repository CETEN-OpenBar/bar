package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var addr = flag.String("addr", "localhost:3737", "http service address")

var listeners = make(map[*websocket.Conn]bool)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func cardReader(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	listeners[c] = true
	defer c.Close()
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
	delete(listeners, c)
}

func server() {
	http.HandleFunc("/", cardReader)
	http.HandleFunc("/config.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "config.json")
	})
	logrus.Fatal(http.ListenAndServe(*addr, nil))
}

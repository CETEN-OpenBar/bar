package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/peterhellberg/acr122u"
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

func main() {
	ctx, err := acr122u.EstablishContext()
	if err != nil {
		panic(err)
	}
	go wsServer()
	h := &handler{logrus.StandardLogger()}
	ctx.Serve(h)
}

func wsServer() {
	http.HandleFunc("/", cardReader)
	logrus.Fatal(http.ListenAndServe(*addr, nil))
}

type handler struct {
	acr122u.Logger
}

func (h *handler) ServeCard(c acr122u.Card) {
	uid := fmt.Sprintf("%x", c.UID())
	logrus.WithField("uid", uid).Info("card inserted")
	for l := range listeners {
		l.WriteJSON(map[string]interface{}{
			"uid": uid,
		})
	}
}

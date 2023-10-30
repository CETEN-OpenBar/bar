package main

import (
	"fmt"

	"github.com/peterhellberg/acr122u"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, err := acr122u.EstablishContext()
	if err != nil {
		panic(err)
	}
	go wsServer()
	h := &handler{logrus.StandardLogger()}
	ctx.Serve(h)
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

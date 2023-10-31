package api

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WSRoom struct {
	conns []*websocket.Conn
}

var wsRooms = make(map[string]*WSRoom)
var lock = sync.RWMutex{}

func GetWSRoom(name string) *WSRoom {
	if _, ok := wsRooms[name]; !ok {
		lock.Lock()
		wsRooms[name] = &WSRoom{}
		lock.Unlock()
	}
	return wsRooms[name]
}

func (s *WSRoom) Add(conn *websocket.Conn) {
	r := *s
	r.conns = append(r.conns, conn)
	*s = r
}

func (s *WSRoom) Remove(conn *websocket.Conn) {
	r := *s
	for i, c := range r.conns {
		if c == conn {
			r.conns = append(r.conns[:i], r.conns[i+1:]...)
			break
		}
	}
	*s = r
}

func (s *WSRoom) Broadcast(message []byte) {
	r := *s
	for _, conn := range r.conns {
		conn.WriteMessage(websocket.TextMessage, message)
	}
	*s = r
}

func BroadcastToRoom(room string, message []byte) {
	r := GetWSRoom(room)
	r.Broadcast(message)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func LinkUpgrade(c echo.Context) error {
	account, err := MustGetUserOrOnBoard(c)
	if err != nil {
		return nil
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return Error500(c)
	}

	room := GetWSRoom(account.Id.String())
	room.Add(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	room.Remove(conn)

	return nil
}

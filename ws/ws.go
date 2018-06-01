package ws

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"fmt"
)

type connection struct {
	conn *websocket.Conn
	send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	c := &connection{
		send: make(chan []byte, 256),
		conn: conn,
	}

	fmt.Println(c)
}

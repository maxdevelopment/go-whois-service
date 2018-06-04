package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"time"
)

const (
	writeWait  = 5 * time.Second
	pingPeriod = 5 * time.Second
)

type Client struct {
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

	client := &Client{
		conn: conn,
		send: make(chan []byte),
	}

	H.register <- client
	go client.listenHub()
	go client.isConnected()
}

func (c *Client) listenHub() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		H.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			for i := 0; i < len(c.send); i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) isConnected() {

	defer func() {
		H.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
	}
}
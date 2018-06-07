package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"time"
	"github.com/gorilla/mux"
	"fmt"
)

const (
	writeWait  = 5 * time.Second
	pingPeriod = 5 * time.Second
)

type Client struct {
	id         string
	conn       *websocket.Conn
	send       chan []byte
	remoteAddr string
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

	params := mux.Vars(r)

	client := &Client{
		id:   params["id"],
		conn: conn,
		send: make(chan []byte),
		//remoteAddr: r.RemoteAddr,
		remoteAddr: "5.61.45.181:5525",
	}

	WhoIs.getData(client)

	H.register <- client
	//service.ConnectedClients <- client
	go client.listenHub()
	//client.isConnected()
}

func (c *Client) listenHub() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		fmt.Println("Unreg deff")
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

//func (c *Client) isConnected() {
//
//	defer func() {
//		fmt.Println("Connected deff")
//		H.unregister <- c
//		c.conn.Close()
//	}()
//
//	for {
//		_, _, err := c.conn.ReadMessage()
//		if err != nil {
//			return
//		}
//	}
//}

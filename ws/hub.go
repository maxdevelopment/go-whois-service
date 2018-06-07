package ws

import (
	"net"
)

type hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
}

var H = hub{
	clients:    make(map[string]*Client),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			_, _, err := net.SplitHostPort(client.remoteAddr)
			if err != nil {
				continue
			}

			h.clients[client.id] = client

			//h.clients[client] = true
			//cc := &service.ConnectedClient{
			//	Ip: client.remoteAddr,
			//	Id: client.id,
			//}
			//service.ConnectedClients <- cc
			//service.ConnectedClients <- client

			//case client := <-h.unregister:
			//	//if _, ok := h.clients[client]; ok {
			//	//	delete(h.clients, client)
			//	//	close(client.send)
			//	//}
			//	//service.DisconnectedClients <- client.id
			//	//service.DisconnectedClients <- client
			//
			//	//case message := <-service.Broadcast:
			//	//	for client := range h.clients {
			//	//		client.send <- message
			//	//	}
		}
	}
}

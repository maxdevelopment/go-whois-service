package ws

import "github.com/maxdevelopment/go-whois-service/service"

type hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

var H = hub{
	clients:    make(map[*Client]bool),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			service.WH.ClientIPs <- client.conn.RemoteAddr()

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}

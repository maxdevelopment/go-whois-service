package ws

import (
	"encoding/json"
)

type hub struct {
	Clients    map[string]*Client `json:"clients"`
	register   chan *Client
	unregister chan *Client
}

var H = hub{
	Clients:    make(map[string]*Client),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func (h *hub) broadcast() {
	clients, _ := json.Marshal(h.Clients)
	for _, client := range h.Clients {
		client.send <- clients
	}
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.Clients[client.Id] = client
			h.broadcast()

		case client := <-h.unregister:
			if _, ok := h.Clients[client.Id]; ok {
				delete(h.Clients, client.Id)
				close(client.send)
			}
			h.broadcast()
		}
	}
}

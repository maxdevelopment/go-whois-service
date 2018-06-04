package ws

import "fmt"

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
			fmt.Println("register", client)
			h.clients[client] = true
			fmt.Println(h.clients)
		case client := <-h.unregister:
			fmt.Println("unregister", client)
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			fmt.Println(h.clients)
		}
	}
}

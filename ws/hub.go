package ws

import "fmt"

type message struct {
	data []byte
	room string
}

type subscription struct {
	conn *connection
	//room string
}

type hub struct {
	rooms map[string]map[*connection]bool
	broadcast chan message
	register chan subscription
	unregister chan subscription
}

var H = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[*connection]bool),
}

func (h *hub) Run() {
	for {
		select {
		case s := <-h.register:
			fmt.Println(s)
		case s := <-h.unregister:
			fmt.Println(s)
		case m := <-h.broadcast:
			fmt.Println(m)
		}
	}
}

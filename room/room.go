package room

import (
	"fmt"
)

const (
	Capacity = 2
)

type room struct {
	msg      chan []byte
	join     chan *client
	leave    chan *client
	clients  map[*client]bool
	capacity int
}

func NewRoom() *room {
	return &room{
		msg:      make(chan []byte),
		join:     make(chan *client),
		leave:    make(chan *client),
		clients:  make(map[*client]bool),
		capacity: Capacity,
	}
}

func (room *room) Run() {
	for {
		select {
		case client := <-room.join:
			room.clients[client] = true
			for c := range room.clients {
				c.send <- []byte(client.name + " join")
			}
			fmt.Println(client.name, ": join room")
		case client := <-room.leave:
			delete(room.clients, client)
			for c := range room.clients {
				c.send <- []byte(client.name + " leave")
			}
			fmt.Println(client.name, ": leave room")
		case msg := <-room.msg:
			for client := range room.clients {
				client.send <- msg
			}
		}
	}
}

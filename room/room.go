package room

import (
	"fmt"
	"net"
	"os"
)

type room struct {
	msg     chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func NewRoom() *room {
	return &room{
		msg:     make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
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

func (room *room) AcceptClient() {
	listener, err := net.Listen("tcp", "localhost:8001")
	defer listener.Close()

	if err != nil {
		println("error listening: ", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			println("error accept: ", err.Error())
			return
		}

		name := make([]byte, 64)
		_, err = conn.Read(name)

		if err != nil {
			fmt.Println("error reading name: ", err.Error())
			return
		}

		client := NewClient(conn, room, string(name))

		room.join <- client

		go client.ReceiveMessage()
		go client.SendMessage()
	}
}

package room

import (
	"fmt"
	"net"
	"os"
)

const (
	RecvBufLen = 1024
)

type client struct {
	conn net.Conn
	name string
	send chan []byte
	room *room
}

func newClient(conn net.Conn, room *room, name string) *client {
	return &client{
		conn: conn,
		name: name,
		send: make(chan []byte),
		room: room,
	}
}

func AcceptClient() {
	listener, err := net.Listen("tcp", "localhost:8001")
	defer listener.Close()

	if err != nil {
		println("error listening: ", err.Error())
		os.Exit(1)
	}

	var rooms []*room

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

		var room *room

		for _, v := range rooms {
			if len(v.clients) < v.capacity {
				room = v
			}
		}

		if room == nil {
			room = NewRoom()
			go room.Run()
			rooms = append(rooms, room)
		}

		client := newClient(conn, room, string(name))

		room.join <- client

		go client.ReceiveMessage()
		go client.SendMessage()
	}
}

func (client *client) ReceiveMessage() {
	defer client.conn.Close()
	for {
		buf := make([]byte, RecvBufLen)
		n, err := client.conn.Read(buf)

		if err != nil {
			fmt.Println("error reading: ", err.Error())
			client.room.leave <- client
			return
		}

		fmt.Println("received ", n, "bytes of data = ", string(buf))

		client.room.msg <- buf
	}
}

func (client *client) SendMessage() {
	for {
		msg := <-client.send
		_, err := client.conn.Write(msg)

		if err != nil {
			fmt.Println("error send reply:", err.Error())
			client.room.leave <- client
			return
		} else {
			fmt.Println("reply sent")
		}

	}
}

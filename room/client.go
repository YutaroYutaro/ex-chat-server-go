package room

import (
	"fmt"
	"net"
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

func NewClient(conn net.Conn, room *room, name string) *client {
	return &client{
		conn: conn,
		name: name,
		send: make(chan []byte),
		room: room,
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

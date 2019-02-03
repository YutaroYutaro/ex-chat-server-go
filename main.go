package main

import "github.com/YutaroYutaro/ex-chat-server-go/room"

func main() {
	println("Run Server.")

	r := room.NewRoom()
	go r.Run()
	r.AcceptClient()
}

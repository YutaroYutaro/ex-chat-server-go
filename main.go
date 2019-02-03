package main

import (
	"fmt"
	"net"
	"os"
)

const (
	RecvBufLen = 1024
)

func main() {
	println("Run Server.")

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

		go EchoFunc(conn)
	}
}

func EchoFunc(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, RecvBufLen)
		n, err := conn.Read(buf)

		if err != nil {
			fmt.Println("error reading: ", err.Error())
			fmt.Println("connect close")
			return
		}

		fmt.Println("received ", n, "bytes of data = ", string(buf))

		_, err = conn.Write(buf)

		if err != nil {
			fmt.Println("error send reply:", err.Error())
		} else {
			fmt.Println("reply sent")
		}
	}
}

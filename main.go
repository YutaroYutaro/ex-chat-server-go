package main

import (
	"net"
	"os"
)

const (
	RecvBufLen = 1024
)

func main() {
	println("Run Server.")

	listener, err := net.Listen("tcp", "localhost:8001")

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
	for {
		buf := make([]byte, RecvBufLen)
		n, err := conn.Read(buf)

		if err != nil {
			println("error reading: ", err.Error())
			return
		}

		println("received ", n, "bytes of data = ", string(buf))

		_, err = conn.Write(buf)

		if err != nil {
			println("error send reply:", err.Error())
		} else {
			println("Reply sent")
		}

		if string(buf) == "exit" {
			conn.Close()
			println("Connect close")
			break
		}
	}
}

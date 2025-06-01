package main

import (
	"fmt"
	"net"
)

func main() {
	h := NewHub()
	go h.Run()

	listener, err := net.Listen("tcp", ":8888")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on port 8888")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
			return
		}
		go HandleConnection(conn, h)
	}
}

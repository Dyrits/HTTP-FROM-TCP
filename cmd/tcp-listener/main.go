package main

import (
	"fmt"
	"github.com/Dyrits/HTTP-FROM-TCP/internal/request"
	"log"
	"net"
)

const file = "messages.txt"

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("An error occurred while listening", err)
	}
	fmt.Println("Listening on", listener.Addr())

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal("An error occurred while closing the listener", err)
		}
		fmt.Println("Listener closed")
	}(listener)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("An error occurred while accepting the connection", err)
		}
		fmt.Println("Accepted connection from:", connection.RemoteAddr())
		request, err := request.RequestFromReader(connection)
		if err != nil {
			log.Println("An error occurred while reading from the connection", err)
		}
		line := request.RequestLine
		line.Print()
		fmt.Println("Closing connection from:", connection.RemoteAddr())
	}
}

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
		lines := readLines(connection)
		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("Closing connection from:", connection.RemoteAddr())
	}

}

func readLines(reader io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {

		defer func(reader io.ReadCloser) {
			err := reader.Close()
			if err != nil {
				fmt.Printf("error: %s\n", err)
			}
		}(reader)

		defer close(lines)

		content := ""
		for {
			bytes := make([]byte, 8, 8)
			_, err := reader.Read(bytes)
			if err != nil {
				if content != "" {
					lines <- content
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Println("An error occurred while reading from the connection", err)
				return
			}
			parts := strings.Split(string(bytes), "\n")
			if len(parts) == 2 {
				lines <- fmt.Sprintf("%s%s", content, parts[0])
				content = ""
			}
			content += parts[len(parts)-1]
		}
	}()
	return lines
}

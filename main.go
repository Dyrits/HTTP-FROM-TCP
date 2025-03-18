package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const file = "messages.txt"

func main() {
	reader, err := os.Open(file)
	if err != nil {
		log.Fatalf("Could not open %s: %s\n", file, err)
	}

	fmt.Printf("Reading data from %s\n", file)

	lines := readLines(reader)

	for line := range lines {
		fmt.Println("read:", line)
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
				fmt.Printf("error: %s\n", err.Error())
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

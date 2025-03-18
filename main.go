package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	var line string

	for {
		bytes := make([]byte, 8)
		_, err := file.Read(bytes)

		if err != nil {
			if err == io.EOF {
				if line != "" {
					fmt.Printf("read: %s\n", line)
				}
				break
			}
			log.Fatal(err)
		}

		parts := strings.Split(string(bytes), "\n")

		if len(parts) == 2 {
			fmt.Printf("read: %s%s\n", line, parts[0])
			line = ""
		}

		line += parts[len(parts)-1]
	}
}

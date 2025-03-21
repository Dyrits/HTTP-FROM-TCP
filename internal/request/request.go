package request

import (
	"fmt"
	"io"
	"log"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var methods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)

	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(data), "\r\n")
	line, err := RequestLineFromString(lines[0])

	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *line,
	}, nil
}

func RequestLineFromString(line string) (*RequestLine, error) {
	if len(line) != 3 {
		return nil, fmt.Errorf("invalid number of parts in request line")
	}

	parts := strings.Split(line, " ")

	method, target, version := parts[0], parts[1], parts[2]

	if !slices.Contains(methods, method) {
		return nil, fmt.Errorf("%s is not a valid HTTP method", method)
	}

	if version != "HTTP/1.1" {
		return nil, fmt.Errorf("%s is not a supported version of HTTP", version)
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}, nil
}

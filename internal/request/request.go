package request

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"
)

type RequestState int

const (
	Initialized RequestState = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	State       RequestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var methods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

func RequestFromReader(reader io.Reader) (*Request, error) {
	// Note: Instead of reading by chunks, we can use a buffered reader, which is more efficient.
	buffer := bufio.NewReader(reader)
	request := &Request{State: Initialized}

	for request.State != Done {
		// Read until a newline is found:
		line, err := buffer.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Parse the line and get the number of bytes parsed:
		parsed, err := request.parse([]byte(line))
		if err != nil {
			return nil, err
		}

		// If more than 0 bytes were parsed, the request is done:
		if parsed > 0 {
			request.State = Done
		}
	}

	if request.State != Done {
		return nil, fmt.Errorf("error: request parsing failed")
	}

	return request, nil
}

func RequestLineFromString(line string) (*RequestLine, error) {
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid number of parts in request line")
	}

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

func (request *Request) parse(data []byte) (int, error) {
	if request.State == Done {
		return 0, fmt.Errorf("error: trying to read data in a done state")
	}

	if request.State != Initialized {
		return 0, fmt.Errorf("error: unknown state")
	}

	// Convert the data to a string and split it into lines:
	lines := strings.SplitN(string(data), "\r\n", 2)
	if len(lines) == 0 {
		return 0, nil
	}

	// Parse the first line as the request line:
	line, err := RequestLineFromString(lines[0])
	if err != nil {
		return 0, err
	}

	// Update the request with the parsed request line:
	request.RequestLine = *line
	request.State = Done

	// Return the number of bytes read:
	return len(lines[0]) + 2, nil
}

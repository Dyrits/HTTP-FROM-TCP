package headers

import (
	"fmt"
	"regexp"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (headers Headers) Parse(data []byte) (number int, done bool, err error) {
	// Look for CRLF:
	CRLF := strings.Index(string(data), "\r\n")

	if CRLF == -1 {
		// No CRLF found, assume more data is needed:
		return 0, false, nil
	}

	if CRLF == 0 {
		// CRLF found at the start, end of headers:
		return 2, true, nil
	}

	// Split the line into key and value, using a regular expression:
	line := string(data)[:CRLF]
	expression := regexp.MustCompile(`^([A-Za-z0-9!#$%&'*+\-.\^_` + "`" + `|~]+):\s*([^\s].*?)\s*$`)
	matches := expression.FindStringSubmatch(line)

	if len(matches) != 3 {
		return 0, false, fmt.Errorf("invalid header format")
	}

	// Extract the key and value:
	key := strings.ToLower(matches[1])
	value := matches[2]

	// Add the key/value pair to the Headers map
	headers[key] = value

	// Return the number of bytes consumed, false for done, and nil for err
	return CRLF + 2, false, nil
}

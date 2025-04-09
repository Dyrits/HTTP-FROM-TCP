package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeadersParser(context *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	number, done, err := headers.Parse(data)
	require.NoError(context, err)
	require.NotNil(context, headers)
	assert.Equal(context, "localhost:42069", headers["host"])
	assert.Equal(context, 23, number)
	assert.False(context, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host:    localhost:42069    \r\n\r\n")
	number, done, err = headers.Parse(data)
	require.NoError(context, err)
	require.NotNil(context, headers)
	assert.Equal(context, "localhost:42069", headers["host"])
	assert.Equal(context, 30, number)
	assert.False(context, done)

	// Test: Valid 2 headers with existing header
	headers = NewHeaders()
	headers["authorization"] = "Bearer <TOKEN>"
	data = []byte("Host: localhost:42069\r\nAccept: application/json\r\n\r\n")
	parsed := 0
	for {
		number, done, err := headers.Parse(data[parsed:])
		require.NoError(context, err)
		parsed += number
		if done {
			break
		}
	}
	require.NoError(context, err)
	require.NotNil(context, headers)
	assert.Equal(context, "localhost:42069", headers["host"])
	assert.Equal(context, "application/json", headers["accept"])
	assert.Equal(context, "Bearer <TOKEN>", headers["authorization"])
	assert.Equal(context, 51, parsed)
	assert.False(context, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n")
	number, done, err = headers.Parse(data)
	require.NoError(context, err)
	assert.Equal(context, 2, number)
	assert.True(context, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	number, done, err = headers.Parse(data)
	require.Error(context, err)
	assert.Equal(context, 0, number)
	assert.False(context, done)

	// Test: Invalid character in header key
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	number, done, err = headers.Parse(data)
	require.Error(context, err)
	assert.Equal(context, 0, number)
	assert.False(context, done)
}

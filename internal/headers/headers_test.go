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
	assert.Equal(context, "localhost:42069", headers["Host"])
	assert.Equal(context, 23, number)
	assert.False(context, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	number, done, err = headers.Parse(data)
	require.Error(context, err)
	assert.Equal(context, 0, number)
	assert.False(context, done)
}

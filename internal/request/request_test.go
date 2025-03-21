package request

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestRequestLineParser(context *testing.T) {
	request, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(context, err)
	require.NotNil(context, request)
	assert.Equal(context, "GET", request.RequestLine.Method)
	assert.Equal(context, "/", request.RequestLine.RequestTarget)
	assert.Equal(context, "HTTP/1.1", request.RequestLine.HttpVersion)

	request, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(context, err)
	require.NotNil(context, request)
	assert.Equal(context, "GET", request.RequestLine.Method)
	assert.Equal(context, "/coffee", request.RequestLine.RequestTarget)
	assert.Equal(context, "HTTP/1.1", request.RequestLine.HttpVersion)

	request, err = RequestFromReader(strings.NewReader("POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(context, err)
	require.NotNil(context, request)
	assert.Equal(context, "POST", request.RequestLine.Method)
	assert.Equal(context, "/coffee", request.RequestLine.RequestTarget)
	assert.Equal(context, "HTTP/1.1", request.RequestLine.HttpVersion)

	_, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(context, err)

	_, err = RequestFromReader(strings.NewReader("/coffee POST HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(context, err)

	_, err = RequestFromReader(strings.NewReader("OPTIONS /prime/rib TCP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(context, err)
}

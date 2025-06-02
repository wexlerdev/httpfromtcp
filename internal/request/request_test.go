package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"strings"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse(t *testing.T) {
	// Test: Good GET Request line
	r, err := RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	r, err = RequestFromReader(strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	//Test: Good POST Request with path
	r, err = RequestFromReader(strings.NewReader("POST /submit/form HTTP/1.1\r\nHost: example.com\r\nContent-Length: 13\r\n\r\nhello=world"))
	require.NoError(t,err)
	require.NotNil(t,r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/submit/form", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	_, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
	require.Error(t, err)

	//Test: Invalid method (out of order)
	_, err = RequestFromReader(strings.NewReader("HTTP/1.1 GET /\r\nHost: example.com\r\n\r\n"))
	require.Error(t, err)

	//Test: Invalid Version in Request Line
	_, err = RequestFromReader(strings.NewReader("GET / HTTP/0.99\r\nHost: example.com\r\n\r\n"))
}


package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"github.com/wexlerdev/httpfromtcp/internal/headers"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	if n > cr.numBytesPerRead {
		n = cr.numBytesPerRead
		cr.pos -= n - cr.numBytesPerRead
	}
	return n, nil
}

func TestRequestLineParse(t *testing.T) {
	// Test: Good GET Request line
	reader := &chunkReader{
		data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead:3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Good GET Request line with path
	reader = &chunkReader {
		data:"GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead:1,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	//Test: Good POST Request with path
	reader = &chunkReader{
		data:"POST /submit/form HTTP/1.1\r\nHost: example.com\r\nContent-Length: 13\r\n\r\nhello=world",
		numBytesPerRead: len("POST /submit/form HTTP/1.1\r\nHost: example.com\r\nContent-Length: 13\r\n\r\nhello=world"),
	}
	r, err = RequestFromReader(reader)
	require.NoError(t,err)
	require.NotNil(t,r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/submit/form", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// Test: Invalid number of parts in request line
	reader = &chunkReader {
		data:"/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 2,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)

	//Test: Invalid method (out of order)
	reader = &chunkReader {
		data: "HTTP/1.1 GET /\r\nHost: example.com\r\n\r\n",
		numBytesPerRead:6,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)

	//Test: Invalid Version in Request Line
	reader = &chunkReader {
		data:"GET / HTTP/0.99\r\nHost: example.com\r\n\r\n",
		numBytesPerRead: 3,
	}
	_, err = RequestFromReader(reader)
}

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	myHeaders := headers.NewHeaders()
	data := []byte("HoSt: localhost:42069\r\n\r\n")
	n, done, err := myHeaders.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, myHeaders)
	assert.Equal(t, "localhost:42069", myHeaders["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	myHeaders = headers.NewHeaders()
	data = []byte("       HOst: localhost:42069                           \r\n\r\n")
	n, done, err = myHeaders.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, myHeaders)
	assert.Equal(t, "localhost:42069", myHeaders["host"])
	assert.Equal(t, 57, n)
	assert.False(t, done)

	// Test: Valid 2 myHeaders with existing myHeaders
	myHeaders = map[string]string{"host": "localhost:42069"}
	data = []byte("User-AgenT: curl/7.81.0\r\nAccept: */*\r\n\r\n")
	n, done, err = myHeaders.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, myHeaders)
	assert.Equal(t, "localhost:42069", myHeaders["host"])
	assert.Equal(t, "curl/7.81.0", myHeaders["user-agent"])
	assert.Equal(t, 25, n)
	assert.False(t, done)

	// Test: Valid done
	myHeaders = headers.NewHeaders()
	data = []byte("\r\n a bunCh of other stuff")
	n, done, err = myHeaders.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, myHeaders)
	assert.Empty(t, myHeaders)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	myHeaders = headers.NewHeaders()
	data = []byte("       HoSt : localhost:42069       \r\n\r\n")
	n, done, err = myHeaders.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	//Test :Invalid chars in field name
	myHeaders = headers.NewHeaders()
	data = []byte("Yams🍠🍠🍠: localhost:69420\r\n\r\n")
	n, done, err = myHeaders.Parse(data)
	require.Error(t, err)
	assert.Equal(t,0,n)
	assert.False(t, done)

	//Test: Valid with multiple values for one field name
	myHeaders = map[string]string{"host": "sillygooses"}
	data = []byte("HosT: moregoosesarehere\r\n\r\n")
	n, done, err = myHeaders.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, myHeaders)
	assert.Equal(t, "sillygooses, moregoosesarehere", myHeaders["host"])
	//assert.Equal(t, 57, n)
	assert.False(t, done)

}

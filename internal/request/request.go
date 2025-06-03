package request

import (
	"io"
	"strings"
	"errors"
	"fmt"
	"github.com/wexlerdev/httpfromtcp/internal/headers"
)
type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	ParserState int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	ParserStateInitialized int = iota // 0
	ParserStateParsingHeaders
	ParserStateDone                   
)

const bufSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufSize, bufSize)
	readToIndex := 0
	req := &Request {
		ParserState:	ParserStateInitialized,
		Headers:		headers.NewHeaders(),
	}

	for req.ParserState != ParserStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf) * 2)
			copy(newBuf, buf)
			buf = newBuf
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.ParserState != ParserStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", req.ParserState, numBytesRead)
				}
				break
			}
		}
		readToIndex += numBytesRead

		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return req, nil
}

func (r * Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.ParserState != ParserStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return totalBytesParsed, nil
		}
		if n == 0 {
			break
		}
		totalBytesParsed += n
	}

	return totalBytesParsed, nil

}

func (r * Request) parseSingle(data []byte) (int, error) {
	totalBytesParsed := 0
	switch r.ParserState {
		case ParserStateDone:
			return 0, errors.New("can't parse in a done state")
		case ParserStateInitialized:
			requestLineStruct, n, err := parseRequestLine(data)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				return 0, nil
			}

			totalBytesParsed += n

			//tis was a success, set the requestLine
			r.RequestLine = * requestLineStruct

			//set state to parsing headers
			r.ParserState = ParserStateParsingHeaders
		case ParserStateParsingHeaders:
			n, done, err  := r.Headers.Parse(data)
			if err != nil {
				return 0, err
			}
			if done {
				r.ParserState = ParserStateDone
			}
			if n == 0 {
				return totalBytesParsed, nil
			}

			totalBytesParsed += n
		default:
			return 0, errors.New("can't parse in an unknown state")
	}
	return totalBytesParsed, nil
}

func parseRequestLine(bites []byte) (*RequestLine, int, error) {
	lineEnd := strings.Index(string(bites), "\r\n")
	if lineEnd == -1 {
		return nil, 0, nil // Not enough data for a complete request line yet
	}

	requestLineStr := string(bites[:lineEnd])
	bytesConsumed := lineEnd + 2 // +2 for "\r\n"

	var requestLineStruct RequestLine
	err := parseRequestLineFromString(requestLineStr, &requestLineStruct)
	if err != nil {
		return nil, 0, err
	}

	return &requestLineStruct, bytesConsumed, nil
}

func parseRequestLineFromString(str string, requestLineStruct *RequestLine) (error) {
	
	sliceOStrings := strings.Split(str, " ")
	if len(sliceOStrings) > 3 {
		return errors.New("Too many parts in the request line, this is wack")
	}

	//method (not method man)
	method := sliceOStrings[0]
	if method != strings.ToUpper(method) {
		return errors.New("method is not all uppercase, this is supa wack")
	}

	//path (feel my wrath on this path)
	path := sliceOStrings[1]
	if path[0] != '/' {
		return errors.New("no leading / in path, I am disapointed in you")
	}

	//http version (I have an aversion to the version, I am certain)
	version := sliceOStrings[2]
	if version != "HTTP/1.1" {
		return errors.New("homie we don't support that version")
	}

	nscanned, err := fmt.Sscanf(version, "HTTP/%s", &version)
	if err != nil {
		return errors.New("whats with the formatting on the version?")
	}
	if nscanned != 1 {
		return errors.New("Something is up with the formatting")
	}
	
	requestLineStruct.Method = method
	requestLineStruct.HttpVersion = version
	requestLineStruct.RequestTarget = path

	return nil
}

package request

import (
	"io"
	"strings"
	"errors"
	"fmt"
)
type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bites, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	requestLine, err := parseRequestLine(bites)
	if err != nil {
		return nil, err
	}

	var requestStruct Request

	requestStruct.RequestLine = *requestLine

	return &requestStruct, nil
	
}

func parseRequestLine(bites []byte) (*RequestLine, error) {
	bigOlString := string(bites)
	sliceOStrings := strings.Split(bigOlString, "\r\n")

	//now we have each line seperate, the request line is at index 0
	requestLine := sliceOStrings[0]

	var requestLineStruct RequestLine
	
	sliceOStrings = strings.Split(requestLine, " ")
	if len(sliceOStrings) > 3 {
		return nil, errors.New("Too many parts in the request line, this is wack")
	}

	//method (not method man)
	method := sliceOStrings[0]
	if method != strings.ToUpper(method) {
		return nil, errors.New("method is not all uppercase, this is supa wack")
	}

	requestLineStruct.Method = method

	//path (feel my wrath on this path)
	path := sliceOStrings[1]
	if path[0] != '/' {
		return nil, errors.New("no leading / in path, I am disapointed in you")
	}

	requestLineStruct.RequestTarget = path

	//http version (I have an aversion to the version, I am certain)
	version := sliceOStrings[2]
	if version != "HTTP/1.1" {
		return nil, errors.New("homie we don't support that version")
	}

	n, err := fmt.Sscanf(version, "HTTP/%s", &version)
	if err != nil {
		return nil, errors.New("whats with the formatting on the version?")
	}
	if n != 1 {
		return nil, errors.New("Something is up with the formatting")
	}
	
	requestLineStruct.HttpVersion = version

	return &requestLineStruct, nil
}

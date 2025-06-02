package request

import (
	"io"
	"strings"
	"errors"
	"fmt"
)
type Request struct {
	RequestLine RequestLine
	ParserState int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	ParserStateInitialized int = iota // 0
	ParserStateDone                   // 1
)

const bufSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {

	var requestStruct Request
	requestStruct.ParserState = ParserStateInitialized

	buf := make([]byte, bufSize)
	readToIndex := 0

	//doesn't handle EOF right now
	for requestStruct.ParserState != ParserStateDone {
		bytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				requestStruct.ParserState = ParserStateDone
				break
			}

			return nil, err
		}
		readToIndex += bytesRead

		bytesParsed, err := requestStruct.parse(buf)
		if err != nil {
			return nil, err
		}
		if bytesParsed == 0 {
			newLen := 2 * cap(buf)

			tempBuf := make([]byte, newLen)
			copy(tempBuf, buf)
			buf = tempBuf
		} else {
			//remove parsed text from buffer
			newLen := cap(buf) - bytesParsed
			if newLen < bufSize {
				newLen = bufSize
			}
			copy(buf, buf[bytesParsed:])
			buf = buf[:newLen]
			readToIndex -= bytesParsed
		}

	}

	return &requestStruct, nil
	
}

func (r * Request) parse(data []byte) (int, error) {
	if r.ParserState == ParserStateDone {
		return 0, errors.New("can't parse in a done state")
	}

	if r.ParserState != ParserStateInitialized {
		return 0, errors.New("can't parse in an unknown state")
	}

	requestLineStruct, n, err := parseRequestLine(data)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, nil
	}

	//tis was a success, set the requestLine
	r.RequestLine = * requestLineStruct

	//set state to done
	r.ParserState = ParserStateDone

	return n, nil

}

func parseRequestLine(bites []byte) (*RequestLine, int, error) {
	bigOlString := string(bites)
	sliceOStrings := strings.Split(bigOlString, "\r\n")

	if len(sliceOStrings) < 2 {
		return nil, 0, nil
	}  

	//now we have each line seperate, the request line is at index 0
	requestLine := sliceOStrings[0]

	var requestLineStruct RequestLine
	n, err := parseRequestLineFromString(requestLine, &requestLineStruct)
	if err != nil {
		return nil, 0, err
	}

	return  &requestLineStruct, n, nil
}

func parseRequestLineFromString(str string, requestLineStruct *RequestLine) (int, error) {
	
	sliceOStrings := strings.Split(str, " ")
	if len(sliceOStrings) > 3 {
		return 0, errors.New("Too many parts in the request line, this is wack")
	}

	//method (not method man)
	method := sliceOStrings[0]
	if method != strings.ToUpper(method) {
		return 0, errors.New("method is not all uppercase, this is supa wack")
	}

	//path (feel my wrath on this path)
	path := sliceOStrings[1]
	if path[0] != '/' {
		return 0, errors.New("no leading / in path, I am disapointed in you")
	}

	//http version (I have an aversion to the version, I am certain)
	version := sliceOStrings[2]
	if version != "HTTP/1.1" {
		return 0, errors.New("homie we don't support that version")
	}

	n, err := fmt.Sscanf(version, "HTTP/%s", &version)
	if err != nil {
		return 0, errors.New("whats with the formatting on the version?")
	}
	if n != 1 {
		return 0, errors.New("Something is up with the formatting")
	}
	
	requestLineStruct.Method = method
	requestLineStruct.HttpVersion = version
	requestLineStruct.RequestTarget = path

	return n, nil
}

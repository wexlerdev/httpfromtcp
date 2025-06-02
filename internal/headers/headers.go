package headers

import (
	"strings"
	"errors"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	rawString := string(data)

	if !strings.Contains(rawString, "\r\n") {
		return 0, false, nil
	}

	//split by crlf but only after one!
	lines := strings.SplitAfter(rawString, "\r\n")

	currentLine := lines[0]
	if strings.HasPrefix(currentLine, "\r\n") {
		return 2, true, nil
	}

	colonSlicedStrings := strings.SplitN(currentLine, ":", 2)
	if len(colonSlicedStrings) != 2 {
		return 0, false, errors.New("parsing error regarding ':'")
	}

	//now we have this notation [0] = field-name [1] = anystring
	
	fieldName := colonSlicedStrings[0]
	//cannot have whitespace on the right aka before colon
	if strings.TrimRightFunc(fieldName, unicode.IsSpace) != fieldName {
		return 0, false, errors.New("cannot have space between field name and :")
	}

	fieldName = strings.TrimSpace(fieldName)
	//fieldName reader for header!
	fieldValue := strings.TrimSpace(colonSlicedStrings[1])
	
	h[fieldName] = fieldValue

	return len(currentLine), false, nil

}





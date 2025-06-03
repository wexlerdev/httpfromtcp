package headers

import (
	"strings"
	"errors"
	"unicode"
	"fmt"
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
	

	if !isValidFieldName(fieldName) {
		return 0, false, errors.New("field name contains invalid chars")
	}


	fieldValue := strings.TrimSpace(colonSlicedStrings[1])


	
	h.Set(fieldName, fieldValue)

	return len(currentLine), false, nil

}

func isValidFieldName(str string) bool {
	var allowedSpecialChars = map[rune]bool{
		'!': true,
		'#': true,
		'$': true,
		'%': true,
		'&': true,
		'\'': true, 
		'*': true,
		'+': true,
		'-': true,
		'.': true,
		'^': true,
		'_': true,
		'`': true,
		'|': true,
		'~': true,
	}
	if len(str) < 1 {
		return false
	}
	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !allowedSpecialChars[r] {
			return false
		}
	}
	return true
}


func (h Headers) Set(key, value string) {

	key = strings.ToLower(key)

	_, exists := h[key]

	if exists {
		value = fmt.Sprintf("%s, %s", h[key], value)
	}
	
	h[key] = value
}

func (h Headers) Get(key string) (string, error){
	lowerKey := strings.ToLower(key)
	value, exists := h[lowerKey]
	if !exists {
		return "", fmt.Errorf("key %s: not found", lowerKey)
	}
	return value, nil
}

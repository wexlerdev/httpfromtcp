package main

import (
	"fmt"
	"os"
	"io"
	"errors"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <- chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		defer f.Close()
				
		currentLine := ""
		for {

			buffah := make([]byte, 8, 8)
			nBytesRead, err := f.Read(buffah)
			if err != nil {
				if currentLine != "" {
					ch <- fmt.Sprintf("%s", currentLine)
					currentLine = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
				return
			}

			str := string(buffah[:nBytesRead])
			parts := strings.Split(str, "\n")
			numParts := len(parts)

			for i := 0; i < numParts -1; i++ {
				ch <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}

			currentLine += parts[numParts-1]
		}

	}()
	return ch
}


func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		os.Exit(1)
	}

	channel := getLinesChannel(file)
	for line := range channel {
		fmt.Printf("read: %s\n", line)
	}
}

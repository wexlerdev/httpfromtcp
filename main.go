package main

import (
	"fmt"
	"net"
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

	listener, err := net.Listen("tcp", ":42069")
	defer listener.Close()

	if err != nil {
		fmt.Printf("ERR: %v\n", err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
			return
		}
		fmt.Println("Connection has been accepted!")



		channel := getLinesChannel(conn)
		for line := range channel {
			fmt.Printf("%s", line)
		}
		fmt.Print("\n")
		fmt.Println("Connection has been closed")
	}

}

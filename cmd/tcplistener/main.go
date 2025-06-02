package main

import (
	"fmt"
	"net"
    "github.com/wexlerdev/httpfromtcp/internal/request" // Import the internal package
)


func main() {

	listener, err := net.Listen("tcp", "localhost:42069")

	if err != nil {
		fmt.Printf("ERR: %v\n", err)
		return
	}
	defer listener.Close()
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
			return
		}
		fmt.Println("Connection has been accepted!")


		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("ERR: %v\n", err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %v \n", req.RequestLine.Method)
		fmt.Printf("- Target: %v \n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v \n", req.RequestLine.HttpVersion)

		fmt.Print("\n")
		fmt.Println("Connection has been closed")
	}

}

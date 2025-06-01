package main

import (
	"net"
	"fmt"
	"bufio"
	"os"
)

func main() {
	udpAddress, err := net.ResolveUDPAddr("udp",":42069")
	if err != nil {
		fmt.Printf("%v \n", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddress)
	if err != nil {
		fmt.Printf("%v \n", err)
	}
	defer conn.Close()

	rawStdin := os.Stdin
	reader := bufio.NewReader(rawStdin)

	for {
		fmt.Print(">")

		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("%v \n", err)
		}

		_, err = conn.Write([]byte(str))
		
		if err != nil {
			fmt.Printf("%v \n", err)
		}
	}


}

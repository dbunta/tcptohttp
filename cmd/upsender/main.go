package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println("Error resolving UDP address")
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error resolving UDP address")
		os.Exit(1)
	}
	defer conn.Close()

	// var reader bufio.Reader
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input")
			continue
			// os.Exit(1)
		}

		_, err = conn.Write([]byte(input))
		if err != nil {
			fmt.Println("Error reading input")
			continue
			// os.Exit(1)
		}
	}

}

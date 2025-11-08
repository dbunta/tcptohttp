package main

import (
	"fmt"
	"net"
	"os"

	"github.com/dbunta/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		fmt.Println("Error creating tcp listener")
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection")
			os.Exit(1)
		}

		fmt.Println("Connection accepted")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for k, v := range req.Headers {
			fmt.Printf("- %v: %v\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Printf(string(req.Body))

		// channel := getLinesChannel(conn)
		// for v := range channel {
		// 	// fmt.Printf("read: %s\n", v)
		// 	fmt.Printf("%s\n", v)
		// }

	}

	// f, err := os.Open("messages.txt")
	// if err != nil {
	// 	fmt.Print("Error opening messages.txt")
	// 	os.Exit(1)
	// }

	// channel := getLinesChannel(f)
	// for v := range channel {
	// 	fmt.Printf("read: %s\n", v)
	// }

	defer listener.Close()
	os.Exit(0)
}

// func getLinesChannel2(f io.ReadCloser) <-chan string {
// 	channel := make(chan string)

// 	go func(channel chan string) {
// 		b := make([]byte, 8)
// 		_, err := f.Read(b)
// 		if err != nil {
// 			// fmt.Print("Error reading messages.txt\n")
// 			fmt.Print(err)
// 			os.Exit(1)
// 		}

// 		var currLine string
// 		for err == nil {
// 			currLine += string(b)
// 			clear(b)
// 			if strings.Contains(currLine, "\n") {
// 				lines := strings.Split(currLine, "\n")
// 				var i int
// 				for i = 0; i < len(lines)-1; i++ {
// 					channel <- lines[i]
// 				}
// 				currLine = lines[len(lines)-1]
// 			}
// 			_, err = f.Read(b)
// 		}

// 		fmt.Println("Connection closed")
// 		defer close(channel)
// 		defer f.Close()
// 	}(channel)

// 	return channel
// }

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	lines := make(chan string)
// 	go func() {
// 		defer f.Close()
// 		defer close(lines)
// 		currentLineContents := ""
// 		for {
// 			b := make([]byte, 8, 8)
// 			n, err := f.Read(b)
// 			if err != nil {
// 				if currentLineContents != "" {
// 					lines <- currentLineContents
// 				}
// 				if errors.Is(err, io.EOF) {
// 					break
// 				}
// 				fmt.Printf("error: %s\n", err.Error())
// 				return
// 			}
// 			str := string(b[:n])
// 			parts := strings.Split(str, "\n")
// 			for i := 0; i < len(parts)-1; i++ {
// 				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
// 				currentLineContents = ""
// 			}
// 			currentLineContents += parts[len(parts)-1]
// 		}
// 	}()
// 	return lines
// }

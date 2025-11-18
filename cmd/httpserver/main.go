package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dbunta/httpfromtcp/internal/request"
	"github.com/dbunta/httpfromtcp/internal/response"
	"github.com/dbunta/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	handler := getDefaultHandler()

	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func getDefaultHandler() server.Handler {
	return server.Handler(func(w io.Writer, req *request.Request) *server.HandlerError {
		herr := server.HandlerError{}

		fmt.Printf("TARGET: %s\r\n", req.RequestLine.RequestTarget)
		if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			fmt.Println("========HERE========")
			route := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
			resp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", route))
			if err != nil {
				herr.StatusCode = response.StatusCode500
				herr.Message = err.Error()
				return &herr
			}
			body := make([]byte, 8)
			n := 1
			for n > 0 {
				n, err = resp.Body.Read(body)
				if err != nil {
					herr.StatusCode = response.StatusCode500
					herr.Message = err.Error()
					return &herr
				}
				herr.StatusCode = response.StatusCode(resp.StatusCode)
				fmt.Printf("Status Code: %v\r\n", resp.StatusCode)
				fmt.Printf("Bytes Read from proxy response: %d\r\n", n)
				if n > 0 {
					w.Write(body)
					m := fmt.Sprintf("\r\n%d\r\n", len(body))
					w.Write([]byte(m))
				} else {
					m := fmt.Sprintf("\r\n%d\r\n\r\n", len(body))
					w.Write([]byte(m))
				}
			}
			herr.StatusCode = response.StatusCode200
		} else if req.RequestLine.RequestTarget == "/yourproblem" {
			herr.StatusCode = response.StatusCode400
			herr.Message = "Your problem is not my problem\n"
			herr.Message = `<html>
				<head>
					<title>400 Bad Request</title>
				</head>
				<body>
					<h1>Bad Request</h1>
					<p>Your request honestly kinda sucked.</p>
				</body>
				</html>`
		} else if req.RequestLine.RequestTarget == "/myproblem" {
			herr.StatusCode = response.StatusCode500
			herr.Message = `<html>
				<head>
					<title>500 Internal Server Error</title>
				</head>
				<body>
					<h1>Internal Server Error</h1>
					<p>Okay, you know what? This one is on me.</p>
				</body>
				</html>`
		} else {
			w.Write([]byte(`<html>
				<head>
					<title>200 OK</title>
				</head>
				<body>
					<h1>Success!</h1>
					<p>Your request was an absolute banger.</p>
				</body>
				</html>`))
			herr.StatusCode = response.StatusCode200
		}

		return &herr
	})
}

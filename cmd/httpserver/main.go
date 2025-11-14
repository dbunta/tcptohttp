package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dbunta/httpfromtcp/internal/request"
	"github.com/dbunta/httpfromtcp/internal/response"
	"github.com/dbunta/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	handler := server.Handler(func(w io.Writer, req *request.Request) *server.HandlerError {
		herr := server.HandlerError{}
		if req.RequestLine.RequestTarget == "/yourproblem" {
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
			// herr.Message = "All good, frfr\n"
		}

		return &herr
	})

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

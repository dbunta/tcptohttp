package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"sync/atomic"

	"github.com/dbunta/httpfromtcp/internal/request"
	"github.com/dbunta/httpfromtcp/internal/response"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
	writer   response.Writer
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func WriteError(w io.Writer, error HandlerError) {
	w.Write([]byte(fmt.Sprintf("An error has occurred\r\n Status Code: %v\r\n %v\r\n", error.StatusCode, error.Message)))
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		return nil, err
	}

	server := Server{
		listener: listener,
	}

	go server.listen(handler)

	return &server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen(handler Handler) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			continue
		}
		go s.handle(conn, handler)
	}
}

func (s *Server) handle(conn net.Conn, handler Handler) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		os.Exit(1)
	}

	var buf []byte
	buffer := bytes.NewBuffer(buf)
	s.writer.Writer = conn

	herr := handler(buffer, req)
	if herr.StatusCode == 200 {
		fmt.Println("HERE")
		headers := response.GetDefaultHeaders(buffer.Len())
		s.writer.WriteStatusLine(response.StatusCode200)
		s.writer.WriteHeaders(headers)
		conn.Write(buffer.Bytes())
	} else {
		fmt.Println("HERE2")
		headers := response.GetDefaultHeaders(len(herr.Message))
		s.writer.WriteStatusLine(herr.StatusCode)
		s.writer.WriteHeaders(headers)
		conn.Write([]byte(herr.Message))
	}
}

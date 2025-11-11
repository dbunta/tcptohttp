package server

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
}

// type serverState int

// const (
// 	serverStateOpen serverState = iota
// 	serverStateClosed
// )

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		return nil, err
	}

	server := Server{
		// closed:    serverStateOpen,
		listener: listener,
	}

	go server.listen()

	return &server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		// if s.state != serverStateClosed {
		if err != nil {
			if s.closed.Load() {
				return
			}
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	// s.state = serverStateClosed
	// var request []byte
	// _, err := conn.Read(request)
	// if err != nil {
	// 	fmt.Printf("Error: %v", err)
	// }
	retval := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!"
	conn.Write([]byte(retval))
	// if err != nil {
	// 	fmt.Printf("Error: %v", err)
	// }
	// fmt.Print(retval)
	// s.state = serverStateOpen
}

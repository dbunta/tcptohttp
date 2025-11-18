package server

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/dbunta/httpfromtcp/internal/headers"
	"github.com/dbunta/httpfromtcp/internal/request"
	"github.com/dbunta/httpfromtcp/internal/response"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
	writer   response.Writer
	handler  Handler
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
		handler:  handler,
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

	req, err := request.RequestFromReader(conn)
	if err != nil {
		os.Exit(1)
	}

	s.writer.Writer = conn
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxy(conn, req, s)
	} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/video") {
		video(conn, req, s)
	} else {
		var buf []byte
		buffer := bytes.NewBuffer(buf)

		herr := s.handler(buffer, req)
		if herr.StatusCode == 200 {
			headers := response.GetDefaultHeaders(buffer.Len())
			s.writer.WriteStatusLine(response.StatusCode200)
			s.writer.WriteHeaders(headers)
			conn.Write(buffer.Bytes())
		} else {
			headers := response.GetDefaultHeaders(len(herr.Message))
			s.writer.WriteStatusLine(herr.StatusCode)
			s.writer.WriteHeaders(headers)
			conn.Write([]byte(herr.Message))
		}
	}
}

func video(conn net.Conn, req *request.Request, s *Server) {
	headers := headers.NewHeaders()
	headers["Content-Type"] = "video/mp4"
	s.writer.WriteStatusLine(response.StatusCode200)
	s.writer.WriteHeaders(headers)

	video, _ := os.ReadFile("./assets/vim.mp4")
	s.writer.Writer.Write(video)

}

func proxy(conn net.Conn, req *request.Request, s *Server) {
	fmt.Println("========HERE========")
	route := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	resp, err := http.Get(fmt.Sprintf("https://httpbin.org/%s", route))
	if err != nil {
	}
	// body := make([]byte, 1024)
	var body []byte
	buf := make([]byte, 1024)
	n := 1
	writeBody := false
	for n > 0 {
		n, err = resp.Body.Read(buf)

		// newbuf := make([]byte, len(body))
		// copy(newbuf[:], body[:])
		// body = make([]byte, 1024+len(body))
		// copy(body[:], newbuf[:])
		// body = append(body, buf[:n]...)

		if !writeBody {
			headers := headers.NewHeaders()
			headers["Transfer-Encoding"] = "chunked"
			contentType := resp.Header.Get("Content-Type")
			headers["Content-Type"] = contentType
			headers["Trailer"] = "X-Content-SHA256, X-Content-Length"
			s.writer.WriteStatusLine(response.StatusCode200)
			s.writer.WriteHeaders(headers)
			writeBody = true
		}
		if err != nil {
		}
		if n > 0 {
			m := fmt.Sprintf("%x\r\n", n)
			s.writer.Writer.Write([]byte(m))
			s.writer.Writer.Write(buf[:n])
			body = append(body, buf[:n]...)
			s.writer.Writer.Write([]byte("\r\n"))
		} else {
			m := fmt.Sprintf("%x\r\n", n)
			s.writer.Writer.Write([]byte(m))
		}
	}
	fmt.Printf("%s\r\n\r\n", body)
	sha := sha256.Sum256(body)
	contentLength := len(body)
	fmt.Printf("%x\r\n", sha[:])
	// s.writer.Writer.Write([]byte("\r\n"))
	s.writer.Writer.Write([]byte(fmt.Sprintf("%s: %x\r\n", "X-Content-SHA256", sha[:])))
	s.writer.Writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", "X-Content-Length", strconv.Itoa(contentLength))))
	s.writer.Writer.Write([]byte("\r\n"))
}

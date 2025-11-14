package response

import (
	"fmt"
	"io"

	"github.com/dbunta/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCode200 StatusCode = 200
	StatusCode400 StatusCode = 400
	StatusCode500 StatusCode = 500
)

type Writer struct {
	Writer io.Writer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	status := ""
	switch statusCode {
	case StatusCode200:
		status = "OK"
	case StatusCode400:
		status = "Bad Request"
	case StatusCode500:
		status = "Internal Server Error"
	}
	_, err := w.Writer.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %v\r\n", statusCode, status)))
	fmt.Printf(fmt.Sprintf("HTTP/1.1 %d %v\r\n", statusCode, status))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["Content-Length"] = fmt.Sprintf("%v", contentLen)
	headers["Connection"] = "close"
	headers["Content-Type"] = "text/html"
	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for k, v := range headers {
		fmt.Printf("%v: %v\r\n", k, v)
		_, err := w.Writer.Write([]byte(fmt.Sprintf("%v: %v\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	w.Writer.Write([]byte("\r\n"))
	fmt.Printf("\r\n")
	return nil
}

func (w *Writer) WriteBody(body []byte) error {
	_, err := w.Writer.Write(body)
	if err != nil {
		return err
	}
	return nil
}

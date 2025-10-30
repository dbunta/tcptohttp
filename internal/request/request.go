package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
	Status      int // 0 == initialized, 1 == done
}
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	// data, err := io.ReadAll(reader)
	// var contents []byte
	// _, err := reader.Read(contents)
	// if err != nil {
	// 	return nil, err
	// }

	// bytesRead, requestLineStr := parseRequestLine(contents)
	// fmt.Println(bytesRead)
	// requestLineItems := strings.Split(requestLineStr, " ")

	var request2 Request
	request2.Status = 0
	readToIndex := 0
	const bufferSize = 2
	buffer := make([]byte, bufferSize)
	var buffer2 []byte

	for {
		// if request is "Done", leave loop
		if request2.Status == 1 {
			break
		}

		if readToIndex >= bufferSize {
			// newBufferSize := len(buffer) * 2
			// buffer = make([]byte, newBufferSize)
			// placeholder := make([]byte, newBufferSize)
			// placeholder := append(buffer2, buffer...)
			// buffer2 = make([]byte, newBufferSize)
			// _ = copy(buffer2, placeholder)
			buffer2 = append(buffer2, buffer...)
		}

		bytesRead, err := reader.Read(buffer)
		if err != nil {
			return nil, err
		}
		readToIndex += bytesRead

		bytesParsed, _ := request2.parse(buffer2)

		if bytesParsed > 0 {
			break
		}
	}

	return &request2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	var bytesParsed int
	if r.Status == 1 {
		return 0, errors.New("error: trying to read data in done state")
	}

	bytesParsed, err := parseRequestLine(r, data)
	if err != nil {
		return 0, err
	}
	if bytesParsed > 0 {
		r.Status = 1
	}

	return bytesParsed, nil
}

// func parseRequestLine(line []byte) (*RequestLine, error) {
func parseRequestLine(r *Request, line []byte) (int, error) {
	requestStr := strings.Split(string(line), "\r\n")
	if len(requestStr) == 1 {
		return 0, nil
	}

	requestLine := requestStr[0]
	requestLineItems := strings.Split(requestLine, " ")

	if len(requestLineItems) < 3 {
		return 0, errors.New("invalid request line")
	}

	method := requestLineItems[0]
	for i := 0; i < len(method); i++ {
		if method[i] > 90 || method[i] < 65 {
			return 0, errors.New("invalid method")
		}
	}

	target := requestLineItems[1]

	version := requestLineItems[2]
	versionParsed := strings.Split(version, "/")
	if len(versionParsed) < 2 || versionParsed[0] != "HTTP" || versionParsed[1] != "1.1" {
		return 0, errors.New("invalid http version")
	}

	r.RequestLine.Method = method
	r.RequestLine.RequestTarget = target
	r.RequestLine.HttpVersion = versionParsed[1]

	return len(line), nil
	// return &RequestLine{
	// 	Method:        method,
	// 	RequestTarget: target,
	// 	HttpVersion:   versionParsed[1],
	// }, 0, nil
}

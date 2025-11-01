package headers

import (
	"bytes"
	"errors"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlf_b := []byte(crlf)
	// if data does not contain crlf
	if !bytes.Contains(data, crlf_b) {
		return 0, false, errors.New("data does not contain crlf")
	}

	// if crlf is at beginning of data, return immediately
	if bytes.HasPrefix(data, crlf_b) {
		// if strings.HasPrefix(header, crlf) {
		return 2, true, nil
	}

	// if !strings.Contains(header, ":") {
	if !bytes.Contains(data, []byte(":")) {
		return 0, false, errors.New("data is in incorrect format (missing key value pair)")
	}

	data = bytes.SplitAfterN(data, crlf_b, 2)[0]
	kvp := bytes.SplitAfterN(data, []byte(":"), 2)
	k := kvp[0]
	k = k[:len(k)-1]
	if unicode.IsSpace(rune(k[len(k)-1])) {
		return 0, false, errors.New("no whitespace allowed between colon and header key")
	}
	h[string(bytes.TrimSpace(k))] = string(bytes.TrimSpace(kvp[1]))

	return len(data), false, nil
}

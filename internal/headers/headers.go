package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
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
		return 0, false, nil
		// return 0, false, errors.New("data does not contain crlf")
	}

	// if crlf is at beginning of data, return immediately
	if bytes.HasPrefix(data, crlf_b) {
		return 2, true, nil
	}

	if !bytes.Contains(data, []byte(":")) {
		return 0, false, errors.New("data is in incorrect format (missing key value pair)")
	}

	data = bytes.SplitAfterN(data, crlf_b, 2)[0]
	kvp := bytes.SplitAfterN(data, []byte(":"), 2)
	k := bytes.TrimSpace(kvp[0])

	if len(k) == 0 {
		return 0, false, errors.New("header must have length > 0")
	}

	k = k[:len(k)-1]
	if unicode.IsSpace(rune(k[len(k)-1])) {
		return 0, false, errors.New("no whitespace allowed between colon and header key")
	}

	fmt.Printf("Header Key: %v", string(k))
	for _, c := range k {
		if (c >= 97 && c <= 122) || (c >= 65 && c <= 90) || (c > 47 && c < 58) || c == 33 || (c >= 35 && c <= 39) || c == 42 || c == 43 || c == 45 || c == 46 || (c >= 94 && c <= 96) || c == 124 || c == 126 {
			continue
		}
		return 0, false, fmt.Errorf("header key contains invalid character: %v", string(c))
	}

	k_str := strings.ToLower(string(k))
	val, ok := h[k_str]
	if ok {
		h[k_str] = val + "," + string(bytes.TrimSpace(kvp[1]))
	} else {
		h[k_str] = string(bytes.TrimSpace(kvp[1]))
	}
	// h[strings.ToLower(string(k))] = string(bytes.TrimSpace(kvp[1]))
	return len(data), false, nil
}

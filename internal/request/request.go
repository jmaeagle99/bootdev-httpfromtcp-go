package request

import (
	"bytes"
	"fmt"
	"io"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type RequestPart int

const (
	RequestPartRequest RequestPart = iota
	RequestPartRequestLine
	RequestPartMethod
	RequestPartProtocol
	RequestPartProtocolVersion
)

type RequestParserError struct {
	Part    RequestPart
	Message string
}

func (e RequestParserError) Error() string {
	return e.Message
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	requestData, err := io.ReadAll(reader)
	if err != nil {
		return nil, RequestParserError{
			RequestPartRequest,
			fmt.Sprintf("Failed to read request from reader: %v", err),
		}
	}
	before, _, ok := bytes.Cut(requestData, []byte("\r\n"))
	if !ok {
		return nil, RequestParserError{
			RequestPartRequest,
			"Request line not found.",
		}
	}
	requestLine, err := parseRequestLine(before)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: requestLine,
	}, nil
}

func parseRequestLine(requestLine []byte) (RequestLine, error) {
	requestLineParts := bytes.Split(requestLine, []byte(" "))
	// Validate part count
	if len(requestLineParts) != 3 {
		return RequestLine{}, RequestParserError{
			Part:    RequestPartRequestLine,
			Message: fmt.Sprintf("Request line contains %d parts when expecting 3 parts.", len(requestLineParts)),
		}
	}
	// Validate HTTP method
	for methodIndex, methodChar := range requestLineParts[0] {
		if methodChar < 'A' || methodChar > 'Z' {
			return RequestLine{}, RequestParserError{
				Part:    RequestPartMethod,
				Message: fmt.Sprintf("Character '%c' at index %d of HTTP method is not valid.", methodChar, methodIndex),
			}
		}
	}
	// Validate HTTP protocol and version
	versionPrefix := []byte("HTTP/")
	if !bytes.HasPrefix(requestLineParts[2], versionPrefix) {
		return RequestLine{}, RequestParserError{
			Part:    RequestPartProtocol,
			Message: fmt.Sprintf("Protocol does not start with '%s'", versionPrefix),
		}
	}
	httpVersion := requestLineParts[2][len(versionPrefix):]
	if !bytes.HasPrefix(httpVersion, []byte("1.1")) {
		return RequestLine{}, RequestParserError{
			Part:    RequestPartProtocolVersion,
			Message: fmt.Sprintf("HTTP version %s is not supported.", httpVersion),
		}
	}

	return RequestLine{
		Method:        string(requestLineParts[0]),
		RequestTarget: string(requestLineParts[1]),
		HttpVersion:   string(httpVersion),
	}, nil
}

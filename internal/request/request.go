package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
)

type RequestParserState int

const (
	RequestParserStateInitialized RequestParserState = iota
	RequestParserStateDone
)

type Request struct {
	RequestLine RequestLine
	ParserState RequestParserState
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

func (r *Request) parse(data []byte) (int, error) {
	switch r.ParserState {
	case RequestParserStateDone:
		return 0, fmt.Errorf("Parser is finished.")
	case RequestParserStateInitialized:
		requestLine, byteCount, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if byteCount == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.ParserState = RequestParserStateDone
		return byteCount, nil
	default:
		return 0, fmt.Errorf("Unknown parser state.")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	readerBuffer := make([]byte, 8)
	requestBuffer := make([]byte, 0, 8)
	readToIndex := 0
	request := Request{
		ParserState: RequestParserStateInitialized,
	}
	for request.ParserState != RequestParserStateDone {
		bytesRead, err := reader.Read(readerBuffer)
		if errors.Is(err, io.EOF) {
			request.ParserState = RequestParserStateDone
			continue
		}
		if err != nil {
			return &Request{}, RequestParserError{
				RequestPartRequest,
				fmt.Sprintf("Failed to read request from reader: %v", err),
			}
		}
		readToIndex += bytesRead
		if readToIndex > cap(requestBuffer) {
			requestBuffer = slices.Grow(requestBuffer, cap(requestBuffer))
		}
		requestBuffer = append(requestBuffer, readerBuffer[:bytesRead]...)

		bytesParsed, err := request.parse(requestBuffer)
		if err != nil {
			return &Request{}, err
		}
		if bytesParsed > 0 {
			requestBuffer = requestBuffer[bytesParsed:]
			readToIndex -= bytesParsed
		}
	}
	return &request, nil
}

func parseRequestLine(requestData []byte) (*RequestLine, int, error) {
	requestLineEndIndex := bytes.Index(requestData, []byte("\r\n"))
	if requestLineEndIndex == -1 {
		return nil, 0, nil
	}
	requestLineParts := bytes.Split(requestData[:requestLineEndIndex], []byte(" "))
	// Validate part count
	if len(requestLineParts) != 3 {
		return nil, 0, RequestParserError{
			Part:    RequestPartRequestLine,
			Message: fmt.Sprintf("Request line contains %d parts when expecting 3 parts.", len(requestLineParts)),
		}
	}
	// Validate HTTP method
	for methodIndex, methodChar := range requestLineParts[0] {
		if methodChar < 'A' || methodChar > 'Z' {
			return nil, 0, RequestParserError{
				Part:    RequestPartMethod,
				Message: fmt.Sprintf("Character '%c' at index %d of HTTP method is not valid.", methodChar, methodIndex),
			}
		}
	}
	// Validate HTTP protocol and version
	versionPrefix := []byte("HTTP/")
	if !bytes.HasPrefix(requestLineParts[2], versionPrefix) {
		return nil, 0, RequestParserError{
			Part:    RequestPartProtocol,
			Message: fmt.Sprintf("Protocol does not start with '%s'", versionPrefix),
		}
	}
	httpVersion := requestLineParts[2][len(versionPrefix):]
	if !bytes.HasPrefix(httpVersion, []byte("1.1")) {
		return nil, 0, RequestParserError{
			Part:    RequestPartProtocolVersion,
			Message: fmt.Sprintf("HTTP version %s is not supported.", httpVersion),
		}
	}

	return &RequestLine{
		Method:        string(requestLineParts[0]),
		RequestTarget: string(requestLineParts[1]),
		HttpVersion:   string(httpVersion),
	}, requestLineEndIndex + 2, nil
}

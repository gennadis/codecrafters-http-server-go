package handler

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// RequestLine represents the request line of an HTTP request.
type RequestLine struct {
	Method   string
	Path     string
	Protocol string
}

// HTTPRequest represents an HTTP request.
type HTTPRequest struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        string
}

// parseHTTPRequest parses the incoming HTTP request from the reader.
func parseHTTPRequest(reader *bufio.Reader) (*HTTPRequest, error) {
	rawRequestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading request line: %v", err)
	}

	requestLine, err := parseRequestLine(rawRequestLine)
	if err != nil {
		return nil, fmt.Errorf("error parsing request line: %v", err)
	}

	headers, err := parseHeaders(reader)
	if err != nil {
		return nil, fmt.Errorf("error parsing request headers: %v", err)
	}

	body, err := parseBody(reader, headers)
	if err != nil {
		return nil, fmt.Errorf("error parsing request body: %v", err)
	}

	return &HTTPRequest{
		RequestLine: *requestLine,
		Headers:     headers,
		Body:        body,
	}, nil
}

// parseRequestLine parses the request line of the HTTP request.
func parseRequestLine(rawRequestLine string) (*RequestLine, error) {
	parts := strings.Fields(rawRequestLine)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid request line: %s", rawRequestLine)
	}
	return &RequestLine{
		Method:   parts[0],
		Path:     parts[1],
		Protocol: parts[2]}, nil
}

// parseHeaders parses the headers of the HTTP request.
func parseHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	return headers, nil
}

// parseBody parses the body of the HTTP request based on the Content-Length header.
func parseBody(reader *bufio.Reader, headers map[string]string) (string, error) {
	contentLengthStr, ok := headers[HeaderContentLength]
	if !ok {
		return "", nil
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return "", fmt.Errorf("invalid Content-Length: %v", err)
	}

	body := make([]byte, contentLength)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %v", err)
	}

	return string(body), nil
}

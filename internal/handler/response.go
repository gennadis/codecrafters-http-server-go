package handler

import (
	"fmt"
	"strings"
)

// HTTPHeader represents a single HTTP header.
type HTTPHeader struct {
	Key   string
	Value string
}

// HTTPResponse represents an HTTP response.
type HTTPResponse struct {
	Status  string
	Headers []HTTPHeader
	Body    string
}

// AddHeader adds a header to the HTTP response.
func (resp *HTTPResponse) AddHeader(key, value string) {
	resp.Headers = append(resp.Headers, HTTPHeader{Key: key, Value: value})
}

// Build constructs the HTTP response string.
func (resp *HTTPResponse) Build() string {
	var headersBuilder strings.Builder
	headersBuilder.WriteString(fmt.Sprintf("HTTP/1.1 %s\r\n", resp.Status))
	for _, header := range resp.Headers {
		headersBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", header.Key, header.Value))
	}
	headersBuilder.WriteString("\r\n")
	return headersBuilder.String() + resp.Body
}

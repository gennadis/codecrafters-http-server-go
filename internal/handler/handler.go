package handler

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/internal/utils"
)

// Endpoints
const (
	EndpointRoot        = "/"
	EndpointUserAgent   = "/user-agent"
	EndpointEchoPrefix  = "/echo/"
	EndpointFilesPrefix = "/files/"
)

// HandleConnection handles incoming connections and processes HTTP requests.
func HandleConnection(conn net.Conn) {
	defer conn.Close()

	rd := bufio.NewReader(conn)
	req, err := parseHTTPRequest(rd)
	if err != nil {
		log.Printf("failed to parse incoming HTTP request: %v", err)
		return
	}
	log.Printf("HTTP request accepted: method %s, path %s", req.RequestLine.Method, req.RequestLine.Path)

	resp := route(req)

	if _, err := conn.Write([]byte(resp.Build())); err != nil {
		log.Printf("error writing response: %v", err)
	}
	log.Printf("HTTP response sent: %s", resp.Build())
}

// route routes HTTP request based on the request path.
func route(req *HTTPRequest) *HTTPResponse {
	switch path := req.RequestLine.Path; {
	case path == EndpointRoot:
		return handleRoot()

	case path == EndpointUserAgent:
		return handleUserAgent(req)

	case strings.HasPrefix(path, EndpointEchoPrefix):
		return handleEcho(path)

	case strings.HasPrefix(path, EndpointFilesPrefix):
		return handleFiles(req)

	default:
		return handleNotFound()
	}
}

// handleRoot handles the root endpoint.
func handleRoot() *HTTPResponse {
	resp := &HTTPResponse{
		Status: StatusOK,
	}
	resp.AddHeader(HeaderContentType, ContentTypeTextPlain)
	resp.AddHeader(HeaderContentLength, "0")
	return resp
}

// handleEcho handles the echo endpoint.
func handleEcho(path string) *HTTPResponse {
	echoStr := strings.TrimPrefix(path, EndpointEchoPrefix)
	resp := &HTTPResponse{
		Status: StatusOK,
		Body:   echoStr,
	}
	resp.AddHeader(HeaderContentType, ContentTypeTextPlain)
	resp.AddHeader(HeaderContentLength, fmt.Sprintf("%d", len(echoStr)))
	return resp
}

// handleUserAgent handles the user-agent endpoint.
func handleUserAgent(req *HTTPRequest) *HTTPResponse {
	userAgent := req.Headers[HeaderUserAgent]
	resp := &HTTPResponse{
		Status: StatusOK,
		Body:   userAgent,
	}
	resp.AddHeader(HeaderContentType, ContentTypeTextPlain)
	resp.AddHeader(HeaderContentLength, fmt.Sprintf("%d", len(userAgent)))
	return resp
}

// handleFiles handles the files endpoint.
func handleFiles(req *HTTPRequest) *HTTPResponse {
	switch req.RequestLine.Method {
	case MethodGet:
		return handleFileGet(req)

	case MethodPost:
		return handleFilePost(req)

	default:
		return handleMethodNotAllowed()
	}
}

// handleFileGet handles GET requests for the files endpoint.
func handleFileGet(req *HTTPRequest) *HTTPResponse {
	filename := strings.TrimPrefix(req.RequestLine.Path, EndpointFilesPrefix)
	fileContent, err := utils.ReadFileContent(filename)
	if errors.Is(err, os.ErrNotExist) {
		return handleNotFound()
	} else if err != nil {
		return handleInternalServerError()
	}

	resp := &HTTPResponse{
		Status: StatusOK,
		Body:   string(fileContent),
	}
	resp.AddHeader(HeaderContentType, ContentTypeApplicationOctetStream)
	resp.AddHeader(HeaderContentLength, fmt.Sprintf("%d", len(fileContent)))
	return resp
}

// handleFilePost handles POST requests for the files endpoint.
func handleFilePost(req *HTTPRequest) *HTTPResponse {
	filename := strings.TrimPrefix(req.RequestLine.Path, EndpointFilesPrefix)
	if err := utils.WriteFileContent(filename, req.Body); err != nil {
		return handleInternalServerError()
	}

	resp := &HTTPResponse{
		Status: StatusCreated,
	}
	resp.AddHeader(HeaderContentType, ContentTypeTextPlain)
	resp.AddHeader(HeaderContentLength, "0")
	return resp
}

// handleNotFound handles 404 Not Found responses.
func handleNotFound() *HTTPResponse {
	resp := &HTTPResponse{
		Status: StatusNotFound,
	}
	resp.AddHeader(HeaderContentType, ContentTypeTextPlain)
	resp.AddHeader(HeaderContentLength, "0")
	return resp
}

// handleInternalServerError handles 500 Internal Server Error responses.
func handleInternalServerError() *HTTPResponse {
	resp := &HTTPResponse{
		Status: StatusInternalServerError,
	}
	resp.AddHeader(HeaderContentType, ContentTypeTextPlain)
	resp.AddHeader(HeaderContentLength, "0")
	return resp
}

// handleMethodNotAllowed handles 405 Method Not Allowed responses.
func handleMethodNotAllowed() *HTTPResponse {
	resp := &HTTPResponse{
		Status: StatusMethodNotAllowed,
	}
	resp.AddHeader(HeaderContentType, ContentTypeTextPlain)
	resp.AddHeader(HeaderContentLength, "0")
	return resp
}

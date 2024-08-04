package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	network    = "tcp"
	listenAddr = "0.0.0.0:4221"
)

type RequestLine struct {
	Method   string
	Path     string
	Protocol string
}

type HTTPRequest struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        string
}

func main() {
	l, err := net.Listen(network, listenAddr)
	if err != nil {
		log.Fatalf("failed to bind to port: %v", err)
	}
	log.Printf("app started on %s", listenAddr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	request, err := parseHTTPRequest(reader)
	if err != nil {
		log.Printf("failed to parse incoming HTTP request: %v", err)
		return
	}
	log.Printf("HTTP request accepted: method %s, path %s", request.RequestLine.Method, request.RequestLine.Path)

	response := generateResponse(request)
	if _, err := conn.Write([]byte(response)); err != nil {
		log.Printf("error writing response: %v", err)
	}
	log.Printf("HTTP response sent: %s", response)
}

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

func parseRequestLine(rawRequestLine string) (*RequestLine, error) {
	parts := strings.Fields(rawRequestLine)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid request line: %s", rawRequestLine)
	}
	return &RequestLine{Method: parts[0], Path: parts[1], Protocol: parts[2]}, nil
}

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

func parseBody(reader *bufio.Reader, headers map[string]string) (string, error) {
	contentLengthStr, ok := headers["Content-Length"]
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

func generateResponse(request *HTTPRequest) string {
	switch path := request.RequestLine.Path; {
	case path == "/":
		return "HTTP/1.1 200 OK\r\n\r\n"

	case strings.HasPrefix(path, "/echo/"):
		echoStr := strings.TrimPrefix(path, "/echo/")
		return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echoStr), echoStr)

	case path == "/user-agent":
		userAgent := request.Headers["User-Agent"]
		return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)

	case strings.HasPrefix(path, "/files/"):
		switch request.RequestLine.Method {
		case "GET":
			filename := strings.TrimPrefix(path, "/files/")

			fileContent, err := readFileContent(filename)
			if errors.Is(err, os.ErrNotExist) {
				return "HTTP/1.1 404 Not Found\r\n\r\n"
			} else if err != nil {
				return "HTTP/1.1 500 Internal Server Error\r\n\r\n"
			}
			return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(fileContent), fileContent)

		case "POST":
			filename := strings.TrimPrefix(path, "/files/")

			if err := writeFileContent(filename, request.Body); err != nil {
				return "HTTP/1.1 500 Internal Server Error\r\n\r\n"
			}

			return "HTTP/1.1 201 Created\r\n\r\n"

		default:
			return "HTTP/1.1 405 Method Not Allowed\r\n\r\n"
		}

	default:
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}
}

func readFileContent(filename string) ([]byte, error) {
	dir := os.Args[2]
	filePath := dir + filename

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("failed to read file %s: %v", filePath, err)
		return nil, err
	}

	return fileContent, nil
}

func writeFileContent(filename string, content string) error {
	dir := os.Args[2]
	filePath := dir + filename

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		log.Printf("failed to write file %s: %v", filePath, err)
		return err
	}

	return nil
}

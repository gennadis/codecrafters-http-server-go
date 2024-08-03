package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
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
	ReqLine RequestLine
	Headers map[string]string
	Body    string
}

func main() {
	l, err := net.Listen(network, listenAddr)
	if err != nil {
		log.Fatalf("failed to bind to port")
	}
	log.Printf("app started on %s", listenAddr)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatalf("error accepting connection %v", err)
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	reader := bufio.NewReader(c)
	request, err := parseRequest(reader)
	if err != nil {
		log.Printf("failed to parse incoming HTTP request: %v", err)
	}
	log.Printf("HTTP request accepted: method %s, urlPath %s", request.ReqLine.Method, request.ReqLine.Path)

	var resp string

	switch path := request.ReqLine.Path; {
	case path == "/":
		resp = "HTTP/1.1 200 OK\r\n\r\n"

	case strings.HasPrefix(path, "/echo/"):
		echoStr := strings.TrimPrefix(path, "/echo/")
		resp = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echoStr), echoStr)

	case path == "/user-agent":
		userAgent := request.Headers["User-Agent"]
		resp = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)

	default:
		resp = "HTTP/1.1 404 Not Found\r\n\r\n"
	}

	if _, err := c.Write([]byte(resp)); err != nil {
		log.Fatalf("error writing response %v", err)
	}
	log.Printf("HTTP response sent: %s", resp)
}

func parseRequest(r *bufio.Reader) (*HTTPRequest, error) {
	var httpReq HTTPRequest
	rawReqLine, err := r.ReadString('\n')
	if err != nil {
		log.Printf("error reading request line: %v", err)
		return nil, err
	}

	reqLine, err := readRequestLine(rawReqLine)
	if err != nil {
		log.Printf("error parsing request line: %v", err)
		return nil, err
	}
	httpReq.ReqLine = *reqLine

	headers, err := readHeaders(r)
	if err != nil {
		log.Printf("error parsing request headers: %v", err)
		return nil, err
	}
	httpReq.Headers = headers

	return &httpReq, nil
}

func readRequestLine(rawReqLine string) (*RequestLine, error) {
	parts := strings.Fields(rawReqLine)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid request line: %s", rawReqLine)
	}
	return &RequestLine{Method: parts[0], Path: parts[1], Protocol: parts[2]}, nil
}

func readHeaders(reader *bufio.Reader) (map[string]string, error) {
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

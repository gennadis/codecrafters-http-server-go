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

func main() {
	l, err := net.Listen(network, listenAddr)
	if err != nil {
		log.Fatalf("failed to bind to port")
	}
	log.Printf("app started on %s", listenAddr)

	c, err := l.Accept()
	if err != nil {
		log.Fatalf("error accepting connection %v", err)
	}
	handleConnection(c)
}

func handleConnection(c net.Conn) {
	defer c.Close()

	reader := bufio.NewReader(c)
	reqLine, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("error reading request line: %v", err)
		return
	}
	method, urlPath, protocol, err := parseRequestLine(reqLine)
	if err != nil {
		log.Printf("error parsing request line: %v", err)
		return
	}
	log.Printf("HTTP request accepted: method %s, urlPath %s, protocol %s", method, urlPath, protocol)

	headers, err := readHeaders(reader)
	if err != nil {
		log.Printf("error parsing request headers: %v", err)
		return
	}

	if urlPath == "/" {
		sendResponse(c, "HTTP/1.1 200 OK\r\n\r\n")
	} else if strings.HasPrefix(urlPath, "/echo") {
		echoStr := strings.TrimPrefix(urlPath, "/echo/")
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echoStr), echoStr)
		sendResponse(c, resp)
	} else if urlPath == "/user-agent" {
		userAgent := headers["User-Agent"]
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
		sendResponse(c, resp)
	} else {
		sendResponse(c, "HTTP/1.1 404 Not Found\r\n\r\n")
	}
}

func sendResponse(c net.Conn, response string) {
	if _, err := c.Write([]byte(response)); err != nil {
		log.Fatalf("error writing response %v", err)
	}
	log.Printf("HTTP response sent: %s", response)
}

func parseRequestLine(reqLine string) (string, string, string, error) {
	// Example request line: "GET /qwerty HTTP/1.1".
	parts := strings.Fields(reqLine)
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("invalid request line: %s", reqLine)
	}
	return parts[0], parts[1], parts[2], nil
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

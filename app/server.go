package main

import (
	"bufio"
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

	// Example request line: "GET /qwerty HTTP/1.1".
	parts := strings.Fields(reqLine)
	if len(parts) < 3 {
		log.Printf("invalid request line: %s", reqLine)
		return
	}

	urlPath := parts[1]
	log.Printf("HTTP request accepted: method %s, urlPath %s, protocol %s", parts[0], urlPath, parts[2])

	switch urlPath {
	case "/":
		sendResponse(c, "HTTP/1.1 200 OK\r\n\r\n")
	default:
		sendResponse(c, "HTTP/1.1 404 Not Found\r\n\r\n")
	}
}

func sendResponse(c net.Conn, response string) {
	if _, err := c.Write([]byte(response)); err != nil {
		log.Fatalf("error writing response %v", err)
	}
	log.Printf("HTTP response sent: %s", response)
}

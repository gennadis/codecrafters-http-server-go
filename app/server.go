package main

import (
	"log"
	"net"
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
	log.Printf("HTTP request accepted")

	okResp := "HTTP/1.1 200 OK\r\n\r\n"
	if _, err := c.Write([]byte(okResp)); err != nil {
		log.Fatalf("error writing response %v", err)
	}
	log.Printf("HTTP response 200 OK sent")
}

package main

import (
	"log"
	"net"

	"github.com/codecrafters-io/http-server-starter-go/internal/handler"
)

const (
	defaultNetwork    = "tcp"
	defaultListenAddr = "127.0.0.1:4221"
)

func main() {
	listener, err := net.Listen(defaultNetwork, defaultListenAddr)
	if err != nil {
		log.Fatalf("failed to bind to port: %v", err)
	}
	log.Printf("server started on %s", defaultListenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection: %v", err)
			continue
		}
		go handler.HandleConnection(conn)
	}
}

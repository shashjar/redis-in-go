package main

import (
	"log"
	"net"
	"os"
)

func configureLogger() {
	log.SetFlags(0)
}

func runServer() {
	log.Println("Running server...")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Println("Failed to bind to port 6379:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func main() {
	configureLogger()
	runServer()
}

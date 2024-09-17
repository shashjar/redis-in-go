package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func readCommand(buf []byte, conn net.Conn) {
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading data into buffer: ", err.Error())
		os.Exit(1)
	}
	log.Printf("Read command: \n%s\n", buf)
}

func writeResponse(buf []byte, conn net.Conn) {
	if string(buf) == "PING" {
		_, err := conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Error writing PONG to connection: ", err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unknown command: \n%s\n", buf)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 128)
		readCommand(buf, conn)
		writeResponse(buf, conn)
	}
}

func main() {
	fmt.Println("Running server...")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

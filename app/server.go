package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func readCommand(conn net.Conn) string {
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading data from connection:", err.Error())
		os.Exit(1)
	}
	message = strings.TrimSuffix(message, "\n")
	fmt.Printf("Read command: %s\n", message)
	return message
}

// TODO: right now just uses \n at end --> implement actual Redis protocol
// TODO: instead of hard-coding these as an if-else block, can use a map of string to function
func writeResponse(message string, conn net.Conn) {
	if strings.ToLower(message) == "ping" {
		_, err := fmt.Fprintf(conn, "+%s\r\n", []byte("PONG"))
		if err != nil {
			fmt.Println("Error writing PONG to connection: ", err.Error())
			os.Exit(1)
		}
	} else if strings.ToLower(message[:4]) == "echo" {
		_, err := fmt.Fprintf(conn, "%s\r\n", []byte(message[5:]))
		if err != nil {
			fmt.Println("Error executing ECHO on connection:", err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unknown command: %s\n", message)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Handling connection")
	for {
		message := readCommand(conn)
		writeResponse(message, conn)
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
			fmt.Println("Error accepting connection:", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

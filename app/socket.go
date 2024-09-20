package main

import (
	"log"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("Handling connection")
	for {
		command, err := readCommand(conn)
		log.Println("Read command:", command, err)
		if err != nil {
			log.Println("Error reading command:", err.Error())
			os.Exit(1)
		}

		if len(command) > 0 {
			executeCommand(command, conn)
		}
	}
}

func readCommand(conn net.Conn) ([]string, error) {
	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return []string{}, nil
	}

	buf = buf[:n]

	cmd, err := parseCommand(buf)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

// TODO: instead of hard-coding these as an if-else block, can use a map of string to function
// TODO: probably want to move this function to a different file eventually (specifically for executing the Redis commands once they've been parsed)
func executeCommand(command []string, conn net.Conn) {
	switch strings.ToLower(command[0]) {
	case "ping":
		_, err := conn.Write([]byte(toSimpleString("PONG")))
		if err != nil {
			log.Println("Error writing PONG to connection: ", err.Error())
			os.Exit(1)
		}
	case "echo":
		_, err := conn.Write([]byte(toBulkString(command[1])))
		if err != nil {
			log.Println("Error executing ECHO on connection:", err.Error())
			os.Exit(1)
		}
	case "command": // TODO: assumes this is COMMAND DOCS and returns an empty array to get redis-cli to work
		_, err := conn.Write([]byte("*0\r\n"))
		if err != nil {
			log.Println("Error executing COMMAND DOCS on connection:", err.Error())
			os.Exit(1)
		}
	default:
		log.Printf("Unknown command: %s\n", command)
		_, err := conn.Write([]byte(toSimpleError("ERR unknown command '" + command[0] + "'")))
		if err != nil {
			log.Println("Error writing to connection: ", err.Error())
			os.Exit(1)
		}
	}
}

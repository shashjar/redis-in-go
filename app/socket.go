package main

import (
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("Handling connection")
	for {
		command, err := readCommand(conn)
		log.Println("Read command:", command, err)
		if err != nil {
			log.Println("Error reading command from connection", err.Error())
			return
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

func write(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println("Error writing to connection:", err.Error())
		conn.Close()
	}
}

// TODO: instead of hard-coding these as an if-else block, can use a map of string to function
// TODO: probably want to move this function to a different file eventually (specifically for executing the Redis commands once they've been parsed)
func executeCommand(command []string, conn net.Conn) {
	switch strings.ToLower(command[0]) {
	case "command":
		switch strings.ToLower(command[1]) {
		case "docs":
			commandDocs(conn)
		default:
			unknownCommand(conn, command)
		}
	case "config":
		switch strings.ToLower(command[1]) {
		case "get":
			configGet(conn, command)
		default:
			unknownCommand(conn, command)
		}
	case "ping":
		ping(conn)
	case "echo":
		echo(conn, command)
	case "get":
		get(conn, command)
	case "set":
		set(conn, command)
	default:
		unknownCommand(conn, command)
	}
}

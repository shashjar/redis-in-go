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

func readIntoBuffer(conn net.Conn) (int, []byte, error) {
	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	if err != nil {
		return n, nil, err
	}

	if n == 0 {
		return 0, []byte{}, nil
	}

	buf = buf[:n]

	return n, buf, nil
}

func readCommand(conn net.Conn) ([]string, error) {
	n, buf, err := readIntoBuffer(conn)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return []string{}, nil
	}

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
// TODO: server can currently error out when accessing command[1] if that wasn't provided - maybe create separate functions as handlers for the top-level commands
// and then have those route to other functions based on the command arguments provided
func executeCommand(command []string, conn net.Conn) {
	// TODO: need to propagate commands to all replicas (if this is a master)
	// TODO: if this is a replica, it should not send responses back to the master after receiving commands
	switch strings.ToLower(command[0]) {
	case "command":
		switch strings.ToLower(command[1]) {
		case "docs":
			commandDocs(conn)
		default:
			invalidCommand(conn, command)
		}
	case "config":
		switch strings.ToLower(command[1]) {
		case "get":
			configGet(conn, command)
		default:
			invalidCommand(conn, command)
		}
	case "info":
		switch strings.ToLower(command[1]) {
		case "replication":
			infoReplication(conn)
		default:
			invalidCommand(conn, command)
		}
	case "replconf":
		replconf(conn)
	case "psync":
		psync(conn)
	case "save":
		save(conn)
	case "ping":
		ping(conn)
	case "echo":
		echo(conn, command)
	case "get":
		get(conn, command)
	case "set":
		set(conn, command)
	case "del":
		del(conn, command)
	case "keys":
		keys(conn, command)
	default:
		unknownCommand(conn, command)
	}
}

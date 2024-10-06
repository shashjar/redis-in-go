package commands

import (
	"net"
	"strings"
)

// Handler for a COMMAND command
func commandHandler(conn net.Conn, command []string) {
	if len(command) < 2 {
		invalidCommand(conn, command)
		return
	}

	switch strings.ToLower(command[1]) {
	case "docs":
		commandDocs(conn)
	default:
		invalidCommand(conn, command)
	}
}

// COMMAND DOCS command - returns an empty array to allow redis-cli to initialize properly
func commandDocs(conn net.Conn) {
	alwaysWrite(conn, "*0\r\n")
}

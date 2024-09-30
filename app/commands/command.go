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

// TODO: returns an empty array to get redis-cli to initialize properly
// COMMAND DOCS command
func commandDocs(conn net.Conn) {
	alwaysWrite(conn, "*0\r\n")
}

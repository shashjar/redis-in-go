package commands

import (
	"log"
	"net"
	"strings"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// TODO: instead of hard-coding these as an if-else block, can use a map of string to function
// and then have those route to other functions based on the command arguments provided
func executeCommand(command []string, conn net.Conn) {
	switch strings.ToLower(command[0]) {
	case "command":
		commandHandler(conn, command)
	case "config":
		configHandler(conn, command)
	case "info":
		infoHandler(conn, command)
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

func unknownCommand(conn net.Conn, command []string) {
	log.Printf("Unknown command: %s\n", command)
	write(conn, protocol.ToSimpleError("ERR unknown command '"+command[0]+"'"))
}

func invalidCommand(conn net.Conn, command []string) {
	log.Printf("Invalid usage of command: %s\n", command)
	write(conn, protocol.ToSimpleError("ERR invalid usage of command '"+command[0]+"'"))
}

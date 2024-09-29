package commands

import (
	"log"
	"net"
	"strings"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// TODO: instead of hard-coding these as an if-else block, can use a map of string to function
// TODO: server can currently error out when accessing command[1] if that wasn't provided - maybe create separate functions as handlers for the top-level commands
// and then have those route to other functions based on the command arguments provided
func executeCommand(command []string, conn net.Conn) {
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

func unknownCommand(conn net.Conn, command []string) {
	log.Printf("Unknown command: %s\n", command)
	write(conn, protocol.ToSimpleError("ERR unknown command '"+command[0]+"'"))
}

func invalidCommand(conn net.Conn, command []string) {
	log.Printf("Invalid usage of command: %s\n", command)
	write(conn, protocol.ToSimpleError("ERR invalid usage of command '"+command[0]+"'"))
}

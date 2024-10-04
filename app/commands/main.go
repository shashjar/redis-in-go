package commands

import (
	"log"
	"net"
	"strings"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/replication"
)

// TODO: instead of hard-coding these as an if-else block, can use a map of string to function
// and then have those route to other functions based on the command arguments provided
func executeCommand(command []string, numCommandBytes int, conn net.Conn) {
	switch strings.ToLower(command[0]) {
	case "multi":
		multi(conn)
		return
	case "exec":
		exec(conn)
		return
	}

	transaction, ok := getOpenTransaction(conn)
	if ok {
		queueCommand(transaction, command, numCommandBytes, conn)
		return
	}

	switch strings.ToLower(command[0]) {
	case "command":
		commandHandler(conn, command)
	case "config":
		configHandler(conn, command)
	case "info":
		infoHandler(conn, command)
	case "replconf":
		replconfHandler(conn, command)
	case "psync":
		psync(conn)
	case "wait":
		wait(conn, command)
	case "save":
		save(conn)
	case "ping":
		ping(conn)
	case "echo":
		echo(conn, command)
	case "type":
		typeCommand(conn, command)
	case "get":
		get(conn, command)
	case "set":
		set(conn, command)
	case "del":
		del(conn, command)
	case "keys":
		keys(conn, command)
	case "incr":
		incr(conn, command)
	case "xadd":
		xadd(conn, command)
	case "xrange":
		xrange(conn, command)
	case "xread":
		xread(conn, command)
	default:
		unknownCommand(conn, command)
	}

	replication.UpdateReplicationOffsetOnReplica(numCommandBytes)
}

func unknownCommand(conn net.Conn, command []string) {
	log.Printf("Unknown command: %s\n", command)
	write(conn, protocol.ToSimpleError("ERR unknown command '"+command[0]+"'"))
}

func invalidCommand(conn net.Conn, command []string) {
	log.Printf("Invalid usage of command: %s\n", command)
	write(conn, protocol.ToSimpleError("ERR invalid usage of command '"+command[0]+"'"))
}

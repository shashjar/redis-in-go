package commands

import (
	"log"
	"net"
	"strings"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/replication"
)

func executeCommand(command []string, numCommandBytes int, transactionExecuting bool, conn net.Conn) {
	if !transactionExecuting {
		commandHandled := handleTransaction(command, numCommandBytes, conn)
		if commandHandled {
			return
		}
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
	case "rpush":
		rpush(conn, command)
	default:
		unknownCommand(conn, command)
	}

	replication.UpdateReplicationOffsetOnReplica(numCommandBytes)
}

func handleTransaction(command []string, numCommandBytes int, conn net.Conn) bool {
	switch strings.ToLower(command[0]) {
	case "multi":
		multi(conn)
		return true
	case "exec":
		exec(conn)
		return true
	case "discard":
		discard(conn)
		return true
	}

	transaction, open := getOpenTransaction(conn)
	if open {
		queueCommand(transaction, command, numCommandBytes, conn)
		return true
	}

	return false
}

func unknownCommand(conn net.Conn, command []string) {
	log.Printf("Unknown command: %s\n", command)
	write(conn, protocol.ToSimpleError("ERR unknown command '"+command[0]+"'"))
}

func invalidCommand(conn net.Conn, command []string) {
	log.Printf("Invalid usage of command: %s\n", command)
	write(conn, protocol.ToSimpleError("ERR invalid usage of command '"+command[0]+"'"))
}

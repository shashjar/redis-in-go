package commands

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/pubsub"
	"github.com/shashjar/redis-in-go/app/replication"
)

func executeCommand(command []string, numCommandBytes int, transactionExecuting bool, conn net.Conn) {
	if !transactionExecuting {
		commandHandled := handleTransaction(command, numCommandBytes, conn)
		if commandHandled {
			return
		}
	}

	commandType := strings.ToLower(command[0])

	inSubscribedMode := pubsub.InSubscribedMode(conn)
	if inSubscribedMode {
		_, ok := pubsub.SUBSCRIBED_MODE_ALLOWED_COMMANDS[commandType]
		if !ok {
			write(conn, protocol.ToSimpleError(fmt.Sprintf("ERR can't execute '%s': only SUBSCRIBE / UNSUBSCRIBE / PING are allowed in this context", commandType)))
			return
		}
	}

	switch commandType {
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
		ping(conn, inSubscribedMode)
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
	case "lpush":
		lpush(conn, command)
	case "lrange":
		lrange(conn, command)
	case "llen":
		llen(conn, command)
	case "lpop":
		lpop(conn, command)
	case "blpop":
		blpop(conn, command)
	case "subscribe":
		subscribe(conn, command)
	case "unsubscribe":
		unsubscribe(conn, command)
	case "publish":
		publish(conn, command)
	case "zadd":
		zadd(conn, command)
	case "zrank":
		zrank(conn, command)
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

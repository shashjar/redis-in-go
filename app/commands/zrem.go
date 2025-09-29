package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// ZREM command
func zrem(conn net.Conn, command []string) {
	if len(command) < 3 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'zrem' command"))
		return
	}

	setKey := command[1]
	members := command[2:]

	numRemoved, errorResponse, ok := store.ZRem(setKey, members)
	if !ok {
		write(conn, protocol.ToSimpleError(errorResponse))
		return
	}

	write(conn, protocol.ToInteger(numRemoved))
}

package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// ZCARD command
func zcard(conn net.Conn, command []string) {
	if len(command) != 2 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'zcard' command"))
		return
	}

	setKey := command[1]

	cardinality, errorResponse, ok := store.ZCard(setKey)
	if !ok {
		write(conn, protocol.ToSimpleError(errorResponse))
		return
	}

	write(conn, protocol.ToInteger(cardinality))
}

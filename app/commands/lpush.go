package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// LPUSH command
func lpush(conn net.Conn, command []string) {
	if len(command) < 3 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'lpush' command"))
		return
	}

	listKey := command[1]
	elements := command[2:]

	newListLength, errorResponse, ok := store.LPush(listKey, elements)
	if ok {
		write(conn, protocol.ToInteger(newListLength))
	} else {
		write(conn, protocol.ToSimpleError(errorResponse))
	}
}

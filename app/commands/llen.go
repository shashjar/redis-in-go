package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// LLEN command
func llen(conn net.Conn, command []string) {
	if len(command) != 2 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'llen' command"))
		return
	}

	listKey := command[1]
	listLength, errorResponse, ok := store.LLen(listKey)
	if ok {
		write(conn, protocol.ToInteger(listLength))
	} else {
		write(conn, protocol.ToSimpleError(errorResponse))
	}
}

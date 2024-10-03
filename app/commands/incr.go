package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// INCR command
func incr(conn net.Conn, command []string) {
	if len(command) != 2 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'incr' command"))
		return
	}

	updatedVal, ok := store.Incr(command[1])
	if !ok {
		write(conn, protocol.ToSimpleError("ERR value is not an integer or out of range"))
		return
	}

	write(conn, protocol.ToInteger(updatedVal))
}

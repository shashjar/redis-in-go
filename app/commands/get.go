package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// GET command
func get(conn net.Conn, command []string) {
	if len(command) != 2 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'get' command"))
		return
	}

	val, ok := store.REDIS_STORE.Get(command[1])
	if ok {
		write(conn, protocol.ToBulkString(val))
	} else {
		write(conn, protocol.ToNullBulkString())
	}
}

package commands

import (
	"fmt"
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// ZSCORE command
func zscore(conn net.Conn, command []string) {
	if len(command) != 3 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'zscore' command"))
		return
	}

	setKey := command[1]
	member := command[2]

	score, memberExists, errorResponse, ok := store.ZScore(setKey, member)
	if !ok {
		write(conn, protocol.ToSimpleError(errorResponse))
		return
	}

	if !memberExists {
		write(conn, protocol.ToNullBulkString())
		return
	}

	write(conn, protocol.ToBulkString(fmt.Sprintf("%v", score)))
}

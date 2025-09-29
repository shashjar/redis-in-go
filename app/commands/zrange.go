package commands

import (
	"net"
	"strconv"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// ZRANGE command
func zrange(conn net.Conn, command []string) {
	if len(command) != 4 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'zrange' command"))
		return
	}

	setKey := command[1]
	startIndexStr := command[2]
	stopIndexStr := command[3]

	startIndex, err := strconv.Atoi(startIndexStr)
	if err != nil {
		write(conn, protocol.ToSimpleError("ERR value is not an integer or out of range"))
		return
	}
	stopIndex, err := strconv.Atoi(stopIndexStr)
	if err != nil {
		write(conn, protocol.ToSimpleError("ERR value is not an integer or out of range"))
		return
	}

	members, errorResponse, ok := store.ZRange(setKey, startIndex, stopIndex)
	if !ok {
		write(conn, protocol.ToSimpleError(errorResponse))
		return
	}

	write(conn, protocol.ToArray(members))
}

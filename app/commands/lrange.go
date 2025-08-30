package commands

import (
	"net"
	"strconv"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// LRANGE command
func lrange(conn net.Conn, command []string) {
	if len(command) != 4 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'lrange' command"))
		return
	}

	listKey := command[1]
	startIndexStr := command[2]
	stopIndexStr := command[3]

	startIndex, err := strconv.Atoi(startIndexStr)
	if err != nil || startIndex < 0 {
		write(conn, protocol.ToSimpleError("ERR value is not an integer or out of range"))
		return
	}
	stopIndex, err := strconv.Atoi(stopIndexStr)
	if err != nil || stopIndex < 0 {
		write(conn, protocol.ToSimpleError("ERR value is not an integer or out of range"))
		return
	}

	elements, errorResponse, ok := store.LRange(listKey, startIndex, stopIndex)
	if ok {
		write(conn, protocol.ToArray(elements))
	} else {
		write(conn, protocol.ToSimpleError(errorResponse))
	}
}

package commands

import (
	"net"
	"strconv"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// BLPOP command
func blpop(conn net.Conn, command []string) {
	if len(command) < 3 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'blpop' command"))
		return
	}

	listKeys := command[1 : len(command)-1]
	timeoutSecStr := command[len(command)-1]

	timeoutSec, err := strconv.Atoi(timeoutSecStr)
	if err != nil || timeoutSec < 0 {
		write(conn, protocol.ToSimpleError("ERR timeout is not an integer or out of range"))
		return
	}

	poppedValue, poppedKey, popped, errorResponse, ok := store.BLPop(listKeys, timeoutSec)
	if !ok {
		write(conn, protocol.ToSimpleError(errorResponse))
		return
	}

	if !popped {
		write(conn, protocol.ToNullBulkString())
	} else {
		write(conn, protocol.ToArray([]string{poppedKey, poppedValue}))
	}
}

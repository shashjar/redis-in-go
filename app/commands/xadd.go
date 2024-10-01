package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// XADD command
func xadd(conn net.Conn, command []string) {
	if len(command) < 5 || len(command)%2 != 1 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'xadd' command"))
		return
	}

	streamKey := command[1]
	entryID := command[2]

	var keys []string
	var values []string
	for i := 3; i < len(command); i += 2 {
		keys = append(keys, command[i])
	}
	for i := 4; i < len(command); i += 2 {
		values = append(values, command[i])
	}

	ok, createdEntryID, errorResponse := store.XAdd(streamKey, entryID, keys, values)
	if ok {
		write(conn, protocol.ToBulkString(createdEntryID))
	} else {
		write(conn, protocol.ToSimpleError(errorResponse))
	}
}

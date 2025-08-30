package commands

import (
	"fmt"
	"net"
	"strconv"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// LPOP command
func lpop(conn net.Conn, command []string) {
	if len(command) != 2 && len(command) != 3 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'lpop' command"))
		return
	}

	listKey := command[1]
	popCount := 1
	if len(command) == 3 {
		count, err := strconv.Atoi(command[2])
		if err != nil || count <= 0 {
			write(conn, protocol.ToSimpleError("ERR value is not an integer or out of range"))
			return
		}
		popCount = count
	}

	poppedElements, errorResponse, ok := store.LPop(listKey, popCount)
	if !ok {
		write(conn, protocol.ToSimpleError(errorResponse))
		return
	}

	if len(poppedElements) == 0 {
		write(conn, protocol.ToNullBulkString())
	} else if len(command) == 2 {
		if len(poppedElements) != 1 {
			panic(fmt.Sprintf("Expected 1 popped element, got %d", len(poppedElements)))
		}
		write(conn, protocol.ToBulkString(poppedElements[0]))
	} else {
		write(conn, protocol.ToArray(poppedElements))
	}
}

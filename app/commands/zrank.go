package commands

import (
	"net"
	"strings"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// ZRANK command
func zrank(conn net.Conn, command []string) {
	withScore := false
	if len(command) < 3 || len(command) > 4 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'zrank' command"))
		return
	} else if len(command) == 4 {
		if strings.ToLower(command[3]) != "withscore" {
			write(conn, protocol.ToSimpleError("ERR syntax error"))
			return
		}
		withScore = true
	}

	setKey := command[1]
	member := command[2]

	rank, score, memberExists, errorResponse, ok := store.ZRank(setKey, member)
	if !ok {
		write(conn, protocol.ToSimpleError(errorResponse))
		return
	}

	if !memberExists {
		write(conn, protocol.ToNullBulkString())
		return
	}

	if withScore {
		write(conn, protocol.ToMixedArray([]interface{}{rank, score}))
	} else {
		write(conn, protocol.ToInteger(rank))
	}
}

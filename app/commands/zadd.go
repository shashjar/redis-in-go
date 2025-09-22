package commands

import (
	"net"
	"strconv"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// ZADD command
func zadd(conn net.Conn, command []string) {
	if len(command) < 4 || len(command)%2 != 0 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'zadd' command"))
		return
	}

	setKey := command[1]
	memberScores := make(map[string]float64)
	for i := 2; i < len(command); i += 2 {
		score, err := strconv.ParseFloat(command[i], 64)
		if err != nil {
			write(conn, protocol.ToSimpleError("ERR value is not a valid float"))
			return
		}
		memberScores[command[i+1]] = score
	}

	numNewMembers, errorResponse, ok := store.ZAdd(setKey, memberScores)
	if !ok {
		write(conn, protocol.ToSimpleError(errorResponse))
		return
	}

	write(conn, protocol.ToInteger(numNewMembers))
}

package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// ECHO command
func echo(conn net.Conn, command []string) {
	if len(command) <= 1 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'echo' command"))
		return
	}

	write(conn, protocol.ToBulkString(command[1]))
}

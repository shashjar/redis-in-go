package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// EXEC command
func exec(conn net.Conn) {
	write(conn, protocol.ToSimpleError("ERR EXEC without MULTI"))
}

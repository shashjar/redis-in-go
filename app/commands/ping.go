package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// PING command
func ping(conn net.Conn) {
	write(conn, protocol.ToSimpleString("PONG"))
}

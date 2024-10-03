package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// MULTI command
func multi(conn net.Conn) {
	write(conn, protocol.ToSimpleString("OK"))
}

package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// REPLCONF command
func replconf(conn net.Conn) {
	write(conn, protocol.ToSimpleString("OK"))
}

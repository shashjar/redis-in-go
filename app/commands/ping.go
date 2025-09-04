package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// PING command
func ping(conn net.Conn, inSubscribedMode bool) {
	if inSubscribedMode {
		write(conn, protocol.ToArray([]string{"PONG", ""}))
	} else {
		write(conn, protocol.ToSimpleString("PONG"))
	}
}

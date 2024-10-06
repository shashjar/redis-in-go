package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// DISCARD command
func discard(conn net.Conn) {
	connectionID := conn.RemoteAddr().String()

	_, open := getOpenTransaction(conn)
	if !open {
		alwaysWrite(conn, protocol.ToSimpleError("ERR DISCARD without MULTI"))
		return
	}

	delete(ACTIVE_TRANSACTIONS, connectionID)
	alwaysWrite(conn, protocol.ToSimpleString("OK"))
}

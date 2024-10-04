package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// EXEC command
func exec(conn net.Conn) {
	connectionID := conn.RemoteAddr().String()

	_, ok := getOpenTransaction(conn)
	if !ok {
		write(conn, protocol.ToSimpleError("ERR EXEC without MULTI"))
		return
	}

	// TODO: execute commands in transaction and write the array of responses to the client
	// transaction := *val
	delete(ACTIVE_TRANSACTIONS, connectionID)
	write(conn, protocol.ToArray([]string{}))
}

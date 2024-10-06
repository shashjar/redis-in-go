package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// EXEC command
func exec(conn net.Conn) {
	connectionID := conn.RemoteAddr().String()

	transaction, open := getOpenTransaction(conn)
	if !open {
		alwaysWrite(conn, protocol.ToSimpleError("ERR EXEC without MULTI"))
		return
	}

	response := executeTransaction(transaction, conn)
	delete(ACTIVE_TRANSACTIONS, connectionID)
	alwaysWrite(conn, response)
}

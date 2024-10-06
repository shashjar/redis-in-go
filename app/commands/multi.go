package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// MULTI command
func multi(conn net.Conn) {
	connectionID := conn.RemoteAddr().String()

	_, open := getOpenTransaction(conn)
	if open {
		alwaysWrite(conn, protocol.ToSimpleError("ERR open transaction already exists for this client"))
		return
	}

	transaction := Transaction{commands: [][]string{}, numCommandBytes: []int{}, responses: []string{}}
	ACTIVE_TRANSACTIONS[connectionID] = &transaction
	alwaysWrite(conn, protocol.ToSimpleString("OK"))
}

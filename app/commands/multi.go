package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// MULTI command
func multi(conn net.Conn) {
	connectionID := conn.RemoteAddr().String()

	_, ok := getOpenTransaction(conn)
	if ok {
		write(conn, protocol.ToSimpleError("ERR open transaction already exists for this client"))
		return
	}

	transaction := Transaction{commands: [][]string{}, numCommandBytes: []int{}}
	ACTIVE_TRANSACTIONS[connectionID] = &transaction
	write(conn, protocol.ToSimpleString("OK"))
}

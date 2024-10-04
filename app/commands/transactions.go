package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// TODO: need to call replication.UpdateReplicationOffsetOnReplica for the numCommandBytes for each command
type Transaction struct {
	commands        [][]string
	numCommandBytes []int
}

// Maps connection ID (remote address as a string) to an open transaction (if one exists) for that client
var ACTIVE_TRANSACTIONS = map[string]*Transaction{}

func getOpenTransaction(clientConn net.Conn) (Transaction, bool) {
	connectionID := clientConn.RemoteAddr().String()
	val, ok := ACTIVE_TRANSACTIONS[connectionID]
	if ok {
		return *val, ok
	} else {
		return Transaction{}, ok
	}
}

func queueCommand(transaction Transaction, command []string, numCommandBytes int, conn net.Conn) {
	transaction.commands = append(transaction.commands, command)
	transaction.numCommandBytes = append(transaction.numCommandBytes, numCommandBytes)

	write(conn, protocol.ToSimpleString("QUEUED"))
}

package commands

import (
	"fmt"
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
)

type Transaction struct {
	commands        [][]string
	numCommandBytes []int
	responses       []string
}

// Maps connection ID (remote address as a string) to an open transaction (if one exists) for that client
var ACTIVE_TRANSACTIONS = map[string]*Transaction{}

func getOpenTransaction(clientConn net.Conn) (*Transaction, bool) {
	connectionID := clientConn.RemoteAddr().String()
	val, ok := ACTIVE_TRANSACTIONS[connectionID]
	if ok {
		return val, ok
	} else {
		return &Transaction{}, ok
	}
}

func queueCommand(transaction *Transaction, command []string, numCommandBytes int, conn net.Conn) {
	transaction.commands = append(transaction.commands, command)
	transaction.numCommandBytes = append(transaction.numCommandBytes, numCommandBytes)

	alwaysWrite(conn, protocol.ToSimpleString("QUEUED"))
}

func executeTransaction(transaction *Transaction, conn net.Conn) string {
	for i, command := range transaction.commands {
		executeCommand(command, transaction.numCommandBytes[i], true, conn)
	}
	response := fmt.Sprintf("%s%d\r\n", protocol.ARRAY, len(transaction.responses))
	for _, commandResponse := range transaction.responses {
		response += commandResponse
	}
	return response
}

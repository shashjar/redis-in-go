package commands

import (
	"net"
	"strconv"
	"strings"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/replication"
)

// Handler for a REPLCONF command
func replconfHandler(conn net.Conn, command []string) {
	if len(command) == 1 {
		replconf(conn)
	} else {
		switch strings.ToLower(command[1]) {
		case "listening-port":
			replconf(conn)
		case "capa":
			replconf(conn)
		case "getack":
			replconfGetAck(conn, command)
		default:
			invalidCommand(conn, command)
		}
	}
}

// REPLCONF command
func replconf(conn net.Conn) {
	write(conn, protocol.ToSimpleString("OK"))
}

// REPLCONF GETACK command
func replconfGetAck(conn net.Conn, command []string) {
	if !replication.SERVER_CONFIG.IsReplica {
		invalidCommand(conn, command)
		return
	}

	alwaysWrite(conn, protocol.ToArray([]string{"REPLCONF", "ACK", strconv.Itoa(replication.SERVER_CONFIG.MasterReplicationOffset)}))
}

package commands

import (
	"log"
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
		case "ack":
			replconfAck(conn, command)
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

// REPLCONF ACK command
func replconfAck(conn net.Conn, ack []string) {
	if len(ack) != 3 {
		log.Println("Error parsing REPLCONF ACK from replica")
		return
	}

	numBytesAcknowledged, err := strconv.Atoi(ack[2])
	if err != nil {
		log.Println("Error parsing number of command bytes acknowledged by replica:", err.Error())
		return
	}

	for _, replica := range replication.SERVER_CONFIG.Replicas {
		if conn.RemoteAddr().String() == replica.ID() {
			replica.LastAcknowledgedReplicationOffset = numBytesAcknowledged
			return
		}
	}
}

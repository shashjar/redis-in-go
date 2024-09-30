package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/replication"
)

// Handler for an INFO command
func infoHandler(conn net.Conn, command []string) {
	if len(command) < 2 {
		invalidCommand(conn, command)
		return
	}

	switch strings.ToLower(command[1]) {
	case "replication":
		infoReplication(conn)
	default:
		invalidCommand(conn, command)
	}
}

// INFO REPLICATION command
func infoReplication(conn net.Conn) {
	var replicationInfo string
	if replication.SERVER_CONFIG.IsReplica {
		replicationInfo = "role:slave\n"
	} else {
		replicationInfo = fmt.Sprintf("role:master\nmaster_replid:%s\nmaster_repl_offset:%d\n", replication.SERVER_CONFIG.MasterReplicationID, replication.SERVER_CONFIG.MasterReplicationOffset)
	}

	write(conn, protocol.ToBulkString(replicationInfo))
}

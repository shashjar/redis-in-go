package commands

import (
	"fmt"
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// INFO REPLICATION command
func infoReplication(conn net.Conn) {
	var replicationInfo string
	if store.SERVER_CONFIG.IsReplica {
		replicationInfo = "role:slave\n"
	} else {
		replicationInfo = fmt.Sprintf("role:master\nmaster_replid:%s\nmaster_repl_offset:%d\n", store.SERVER_CONFIG.MasterReplicationID, store.SERVER_CONFIG.MasterReplicationOffset)
	}

	write(conn, protocol.ToBulkString(replicationInfo))
}

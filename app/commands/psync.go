package commands

import (
	"fmt"
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/replication"
	"github.com/shashjar/redis-in-go/app/store"
)

// PSYNC command
func psync(conn net.Conn) {
	response := fmt.Sprintf("FULLRESYNC %s %d", store.SERVER_CONFIG.MasterReplicationID, store.SERVER_CONFIG.MasterReplicationOffset)
	write(conn, protocol.ToSimpleString(response))
	replication.ExecuteFullResync(conn)
}

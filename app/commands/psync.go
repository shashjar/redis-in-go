package commands

import (
	"fmt"
	"net"
	"time"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/replication"
)

// PSYNC command
func psync(conn net.Conn) {
	response := fmt.Sprintf("FULLRESYNC %s %d", replication.SERVER_CONFIG.MasterReplicationID, replication.SERVER_CONFIG.MasterReplicationOffset)
	write(conn, protocol.ToSimpleString(response))
	time.Sleep(500 * time.Millisecond)
	replication.ExecuteFullResync(conn)
}

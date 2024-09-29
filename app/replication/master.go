package replication

import (
	"log"
	"net"

	"github.com/shashjar/redis-in-go/app/persistence"
)

// Executes a full resynchronization by sending an RDB file from the master to the replica on the given connection.
func ExecuteFullResync(conn net.Conn) {
	rdbFileBytes := persistence.GetRDBBytes()
	_, err := conn.Write(rdbFileBytes)
	if err != nil {
		log.Println("Error writing RDB file to replica:", err.Error())
		conn.Close()
	}
}

// Generates a replication ID for the master server.
func generateMasterReplicationID() string {
	return "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
}

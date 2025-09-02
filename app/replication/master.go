package replication

import (
	"log"
	"net"
	"strings"

	"github.com/shashjar/redis-in-go/app/persistence"
	"github.com/shashjar/redis-in-go/app/protocol"
)

// Represents the set of Redis commands that should be propagated from the master server to replicas
var COMMANDS_TO_PROPAGATE = map[string]struct{}{
	"set":   {},
	"del":   {},
	"incr":  {},
	"xadd":  {},
	"rpush": {},
	"lpush": {},
	"lpop":  {},
	"blpop": {},
}

const NUM_GET_ACK_BYTES = 37

// Executes a full resynchronization by sending an RDB file from the master to the replica on the given connection.
func ExecuteFullResync(conn net.Conn) {
	rdbFileBytes := persistence.GetRDBBytes()
	_, err := conn.Write(rdbFileBytes)
	if err != nil {
		log.Println("Error writing RDB file to replica:", err.Error())
		conn.Close()
	}
}

// Propagates the given command from the master server to all existing replicas
func PropagateCommand(commandName string, commandBytes []byte) {
	if !SERVER_CONFIG.IsReplica {
		_, ok := COMMANDS_TO_PROPAGATE[strings.ToLower(commandName)]
		if ok {
			for _, replica := range SERVER_CONFIG.Replicas {
				_, err := replica.Conn.Write(commandBytes)
				if err != nil {
					log.Println("Error propagating command to replica:", err.Error())
				}
			}
			SERVER_CONFIG.MasterReplicationOffset += len(commandBytes)
		}
	}
}

// Sends the REPLCONF GETACK command to the replicas, intending to verify replication offsets
func SendGetAckToReplicas() {
	for _, replica := range SERVER_CONFIG.Replicas {
		_, err := replica.Conn.Write([]byte(protocol.ToArray([]string{"REPLCONF", "GETACK", "*"})))
		if err != nil {
			log.Println("Error sending REPLCONF GETACK command to replica:", err.Error())
		}
	}
}

// Adds the number of REPLCONF GETACK command bytes to the master's replication offset
func AddGetAckBytesToMasterReplicationOffset() {
	SERVER_CONFIG.MasterReplicationOffset += NUM_GET_ACK_BYTES
}

// Generates a replication ID for the master server
func generateMasterReplicationID() string {
	return "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
}

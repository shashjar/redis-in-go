package replication

import "net"

// Initializes the master or replica server's replication configuration.
func InitializeReplication() net.Conn {
	if SERVER_CONFIG.IsReplica {
		return replicaHandshake()
	} else {
		SERVER_CONFIG.MasterReplicationID = generateMasterReplicationID()
		return nil
	}
}

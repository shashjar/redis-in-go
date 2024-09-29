package replication

import "github.com/shashjar/redis-in-go/app/store"

// Initializes the master or replica server's replication configuration.
func InitializeReplication() {
	if store.SERVER_CONFIG.IsReplica {
		replicaHandshake()
	} else {
		store.SERVER_CONFIG.MasterReplicationID = generateMasterReplicationID()
	}
}

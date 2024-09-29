package replication

// Initializes the master or replica server's replication configuration.
func InitializeReplication() {
	if SERVER_CONFIG.IsReplica {
		replicaHandshake()
	} else {
		SERVER_CONFIG.MasterReplicationID = generateMasterReplicationID()
	}
}

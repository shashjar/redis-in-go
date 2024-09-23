package main

func initializeReplication() {
	if !SERVER_CONFIG.isReplica {
		SERVER_CONFIG.masterReplicationID = generateMasterReplicationID()
	}
}

func generateMasterReplicationID() string {
	return "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
}

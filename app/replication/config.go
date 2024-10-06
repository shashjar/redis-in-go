package replication

import (
	"strings"
)

// Stores the configuration of the current running server
var SERVER_CONFIG = ServerConfig{
	Port:                    "6379",
	IsReplica:               false,
	MasterHost:              "",
	MasterPort:              "",
	MasterReplicationID:     "",
	MasterReplicationOffset: 0,
	Replicas:                []*Replica{},
}

// Given the server's running port and master that it is replicating (if applicable), updates the server configuration
func UpdateServerConfig(portPtr *string, replicaOfPtr *string) {
	SERVER_CONFIG.Port = *portPtr
	if len(*replicaOfPtr) > 0 {
		SERVER_CONFIG.IsReplica = true
		replicaOfInfo := strings.Split(*replicaOfPtr, " ")
		SERVER_CONFIG.MasterHost = replicaOfInfo[0]
		SERVER_CONFIG.MasterPort = replicaOfInfo[1]
	}
}

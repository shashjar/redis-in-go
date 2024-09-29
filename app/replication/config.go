package replication

import (
	"net"
	"strings"
)

var SERVER_CONFIG = ServerConfig{
	Port:                    "6379",
	IsReplica:               false,
	MasterHost:              "",
	MasterPort:              "",
	MasterReplicationID:     "",
	MasterReplicationOffset: 0,
	Replicas:                []net.Conn{},
}

func UpdateServerConfig(portPtr *string, replicaOfPtr *string) {
	SERVER_CONFIG.Port = *portPtr
	if len(*replicaOfPtr) > 0 {
		SERVER_CONFIG.IsReplica = true
		replicaOfInfo := strings.Split(*replicaOfPtr, " ")
		SERVER_CONFIG.MasterHost = replicaOfInfo[0]
		SERVER_CONFIG.MasterPort = replicaOfInfo[1]
	}
}

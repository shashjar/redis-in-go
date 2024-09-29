package replication

import "net"

// Represents the configuration of a Redis server
type ServerConfig struct {
	Port                    string
	IsReplica               bool
	MasterHost              string
	MasterPort              string
	MasterReplicationID     string
	MasterReplicationOffset int
	Replicas                []net.Conn
}

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
	Replicas                []*Replica
}

// Represents a replica Redis server
type Replica struct {
	Conn                              net.Conn
	LastAcknowledgedReplicationOffset int
}

func (replica *Replica) ID() string {
	return replica.Conn.RemoteAddr().String()
}

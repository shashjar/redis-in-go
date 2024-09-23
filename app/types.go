package main

// Represents the configuration of a Redis server
type ServerConfig struct {
	port                    string
	isReplica               bool
	replicaOf               string
	masterReplicationID     string
	masterReplicationOffset int
}

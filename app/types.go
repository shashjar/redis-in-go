package main

import "net"

// Represents the configuration of a Redis server
type ServerConfig struct {
	port                    string
	isReplica               bool
	masterHost              string
	masterPort              string
	masterReplicationID     string
	masterReplicationOffset int
	replicas                []net.Conn
}

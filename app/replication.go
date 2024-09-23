package main

import (
	"log"
	"net"
	"os"
	"time"
)

func initializeReplication() {
	if SERVER_CONFIG.isReplica {
		replicaHandshake()
	} else {
		SERVER_CONFIG.masterReplicationID = generateMasterReplicationID()
	}
}

func replicaHandshake() {
	conn, err := net.Dial(NETWORK, SERVER_CONFIG.masterHost+":"+SERVER_CONFIG.masterPort)
	if err != nil {
		log.Println("Failed to connect to master from replica:", err.Error())
		os.Exit(1)
	}

	// TODO: not checking responses from the master to these, just sleeping and assuming these are all OK
	write(conn, toArray([]string{"PING"}))
	time.Sleep(500 * time.Millisecond)
	write(conn, toArray([]string{"REPLCONF", "listening-port", SERVER_CONFIG.port}))
	time.Sleep(500 * time.Millisecond)
	write(conn, toArray([]string{"REPLCONF", "capa", "psync2"}))
	time.Sleep(500 * time.Millisecond)
	write(conn, toArray([]string{"PSYNC", "?", "-1"}))
}

func generateMasterReplicationID() string {
	return "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
}

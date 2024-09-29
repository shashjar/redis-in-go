package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

/** SHARED **/

// Initializes the master or replica server's replication configuration
func initializeReplication() {
	if SERVER_CONFIG.isReplica {
		replicaHandshake()
	} else {
		SERVER_CONFIG.masterReplicationID = generateMasterReplicationID()
	}
}

/** MASTER **/

// Generates a replication ID for the master server
func generateMasterReplicationID() string {
	return "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
}

// Executes a full resynchronization by sending an RDB file from the master to the replica on the given connection
func executeFullResync(conn net.Conn) {
	rdbFileBytes := getRDBBytes()
	_, err := conn.Write(rdbFileBytes)
	if err != nil {
		log.Println("Error writing RDB file to replica:", err.Error())
		conn.Close()
	}
}

/** REPLICA **/

// Sends the replication handshake from the replica to the master
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
	write(conn, toArray([]string{"REPLCONF", "capa", "eof", "capa", "psync2"}))
	time.Sleep(500 * time.Millisecond)
	write(conn, toArray([]string{"PSYNC", "?", "-1"}))

	readRDBFileFromMaster(conn)
}

// Reads an RDB file from the master server onto this replica from the provided connection, and persists
// the RDB file state into the active key-value store
func readRDBFileFromMaster(conn net.Conn) {
	filePath := "/redis-data/replica/incoming_rdb.rdb"
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Error creating file for reading RDB file from master Redis server into:", err.Error())
	}
	defer file.Close()

	_, err = io.Copy(file, conn)
	if err != nil {
		log.Fatal("Error reading RDB file from master onto replica:", err.Error())
	}

	persistFromRDB(filePath)
}

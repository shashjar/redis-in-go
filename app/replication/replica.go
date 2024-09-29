package replication

import (
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/shashjar/redis-in-go/app/persistence"
	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

const NETWORK string = "tcp"

// Sends the replication handshake from the replica to the master.
func replicaHandshake() {
	conn, err := net.Dial(NETWORK, store.SERVER_CONFIG.MasterHost+":"+store.SERVER_CONFIG.MasterPort)
	if err != nil {
		log.Println("Failed to connect to master from replica:", err.Error())
		os.Exit(1)
	}

	// TODO: not checking responses from the master to these, just sleeping and assuming these are all OK
	writeFromReplica(conn, protocol.ToArray([]string{"PING"}))
	time.Sleep(500 * time.Millisecond)
	writeFromReplica(conn, protocol.ToArray([]string{"REPLCONF", "listening-port", store.SERVER_CONFIG.Port}))
	time.Sleep(500 * time.Millisecond)
	writeFromReplica(conn, protocol.ToArray([]string{"REPLCONF", "capa", "eof", "capa", "psync2"}))
	time.Sleep(500 * time.Millisecond)
	writeFromReplica(conn, protocol.ToArray([]string{"PSYNC", "?", "-1"}))

	readRDBFileFromMaster(conn)
}

func writeFromReplica(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println("Error writing to connection:", err.Error())
		conn.Close()
	}
}

// Reads an RDB file from the master server onto this replica from the provided connection, and persists
// the RDB file state into the active key-value store.
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

	persistence.PersistFromRDB(filePath)
}

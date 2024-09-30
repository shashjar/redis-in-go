package replication

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/shashjar/redis-in-go/app/persistence"
	"github.com/shashjar/redis-in-go/app/protocol"
)

const NETWORK string = "tcp"

// Sends the replication handshake from the replica to the master.
func replicaHandshake() {
	conn, err := net.Dial(NETWORK, SERVER_CONFIG.MasterHost+":"+SERVER_CONFIG.MasterPort)
	if err != nil {
		log.Println("Failed to connect to master from replica:", err.Error())
		os.Exit(1)
	}

	pingMasterServer(conn)
	replconfListeningPort(conn)
	replconfCapabilities(conn)
	psync(conn)

	readRDBFileFromMaster(conn)
}

func pingMasterServer(conn net.Conn) {
	writeFromReplica(conn, protocol.ToArray([]string{"PING"}))
	response, err := readResponseFromMaster(conn)
	if err != nil || response[:5] != "+PONG" {
		conn.Close()
		log.Fatal("Failed to PING master server during replica handshake")
	}
}

func replconfListeningPort(conn net.Conn) {
	writeFromReplica(conn, protocol.ToArray([]string{"REPLCONF", "listening-port", SERVER_CONFIG.Port}))
	response, err := readResponseFromMaster(conn)
	if err != nil || response[:3] != "+OK" {
		conn.Close()
		log.Fatal("Failed to send REPLCONF with listening port to master server during replica handshake")
	}
}

func replconfCapabilities(conn net.Conn) {
	writeFromReplica(conn, protocol.ToArray([]string{"REPLCONF", "capa", "eof", "capa", "psync2"}))
	response, err := readResponseFromMaster(conn)
	if err != nil || response[:3] != "+OK" {
		conn.Close()
		log.Fatal("Failed to send REPLCONF with capabilities to master server during replica handshake")
	}
}

func psync(conn net.Conn) {
	writeFromReplica(conn, protocol.ToArray([]string{"PSYNC", "?", "-1"}))
	response, err := readResponseFromMaster(conn)
	if err != nil || response[:11] != "+FULLRESYNC" {
		conn.Close()
		log.Fatal("Failed to PSYNC with master server during replica handshake")
	}
}

// Reads an RDB file from the master server onto this replica from the provided connection, and persists
// the RDB file state into the active key-value store.
func readRDBFileFromMaster(conn net.Conn) {
	filePath := "./redis-data/replica/incoming_rdb.rdb"
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Error creating file for reading RDB file from master Redis server into:", err.Error())
	}
	defer file.Close()

	_, buf, err := readIntoBuffer(conn)
	if err != nil {
		log.Fatal("Error reading RDB file from master Redis server:", err.Error())
	}
	file.Write(buf)

	persistence.PersistFromRDB(filePath)
}

func readResponseFromMaster(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	return response, err
}

func readIntoBuffer(conn net.Conn) (int, []byte, error) {
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return n, nil, err
	}

	if n == 0 {
		return 0, []byte{}, nil
	}

	buf = buf[:n]

	return n, buf, nil
}

func writeFromReplica(conn net.Conn, message string) {
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println("Error writing to connection:", err.Error())
		conn.Close()
	}
}

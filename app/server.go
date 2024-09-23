package main

import (
	"flag"
	"log"
	"net"
	"os"
)

const NETWORK string = "tcp"
const ADDRESS string = "0.0.0.0"

var SERVER_CONFIG = ServerConfig{
	port:                    "6379",
	isReplica:               false,
	replicaOf:               "",
	masterReplicationID:     "",
	masterReplicationOffset: 0,
}

func configureLogger() {
	log.SetFlags(0)
}

func parseCommandLineArguments() {
	portPtr := flag.String("port", "6379", "Port number on which to run Redis server")
	replicaOfPtr := flag.String("replicaof", "", "Indicates whether this Redis server should assume the 'replica' role")
	rdbDirPtr := flag.String("dir", DEFAULT_RDB_DIR, "Directory in which to store RDB file")
	rdbFilenamePtr := flag.String("dbfilename", DEFAULT_RDB_FILENAME, "Filename for RDB file")

	flag.Parse()

	SERVER_CONFIG.port = *portPtr
	if len(*replicaOfPtr) > 0 {
		SERVER_CONFIG.isReplica = true
		SERVER_CONFIG.replicaOf = *replicaOfPtr
	}
	RDB_DIR = *rdbDirPtr
	RDB_FILENAME = *rdbFilenamePtr
}

func runServer() {
	log.Println("Running server...")

	l, err := net.Listen(NETWORK, ADDRESS+":"+SERVER_CONFIG.port)
	if err != nil {
		log.Println("Failed to bind to port 6379:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

func main() {
	configureLogger()
	parseCommandLineArguments()
	initializeReplication()
	persistFromRDB()
	runServer()
}

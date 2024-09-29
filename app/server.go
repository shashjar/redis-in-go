package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/shashjar/redis-in-go/app/commands"
	"github.com/shashjar/redis-in-go/app/persistence"
	"github.com/shashjar/redis-in-go/app/replication"
	"github.com/shashjar/redis-in-go/app/store"
)

const NETWORK string = "tcp"
const ADDRESS string = "0.0.0.0"

func configureLogger() {
	log.SetFlags(0)
}

func parseCommandLineArguments() {
	portPtr := flag.String("port", "6379", "Port number on which to run Redis server")
	replicaOfPtr := flag.String("replicaof", "", "Indicates whether this Redis server should assume the 'replica' role")
	rdbDirPtr := flag.String("dir", persistence.DEFAULT_RDB_DIR, "Directory in which to store RDB file")
	rdbFilenamePtr := flag.String("dbfilename", persistence.DEFAULT_RDB_FILENAME, "Filename for RDB file")

	flag.Parse()

	store.UpdateServerConfig(portPtr, replicaOfPtr)
	persistence.RDB_DIR = *rdbDirPtr
	persistence.RDB_FILENAME = *rdbFilenamePtr
}

func runServer() {
	log.Println("Running server...")

	l, err := net.Listen(NETWORK, ADDRESS+":"+store.SERVER_CONFIG.Port)
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

		go commands.HandleConnection(conn)
	}
}

// TODO: add nicer documentation throughout code (purpose statements?) and update repo README
func main() {
	configureLogger()
	parseCommandLineArguments()
	replication.InitializeReplication()
	persistence.PersistFromRDB("." + persistence.RDB_DIR + "/" + persistence.RDB_FILENAME)
	runServer()
}

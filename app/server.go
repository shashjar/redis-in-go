package main

import (
	"flag"
	"log"
	"net"
	"os"
)

var REDIS_STORE = KeyValueStore{data: make(map[string]KeyValue)}

func configureLogger() {
	log.SetFlags(0)
}

func parseCommandLineArguments() {
	rdbDirPtr := flag.String("dir", DEFAULT_RDB_DIR, "Directory in which to store RDB file")
	rdbFilenamePtr := flag.String("dbfilename", DEFAULT_RDB_FILENAME, "Filename for RDB file")

	flag.Parse()

	RDB_DIR = *rdbDirPtr
	RDB_FILENAME = *rdbFilenamePtr
}

func runServer() {
	log.Println("Running server...")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
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
	runServer()
}

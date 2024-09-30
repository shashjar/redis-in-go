package commands

import (
	"net"
	"strings"

	"github.com/shashjar/redis-in-go/app/persistence"
	"github.com/shashjar/redis-in-go/app/protocol"
)

// Handler for a CONFIG command
func configHandler(conn net.Conn, command []string) {
	if len(command) < 2 {
		invalidCommand(conn, command)
	}

	switch strings.ToLower(command[1]) {
	case "get":
		configGet(conn, command)
	default:
		invalidCommand(conn, command)
	}
}

// CONFIG GET command
func configGet(conn net.Conn, command []string) {
	if len(command) <= 2 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'config get' command"))
		return
	}

	var configParams []string
	for i := 2; i < len(command); i++ {
		switch command[i] {
		case "dir":
			configParams = append(configParams, "dir")
			configParams = append(configParams, persistence.RDB_DIR)
		case "dbfilename":
			configParams = append(configParams, "dbfilename")
			configParams = append(configParams, persistence.RDB_FILENAME)
		default:
			write(conn, protocol.ToSimpleError("ERR invalid configuration parameter for 'config get' command"))
			return
		}
	}

	write(conn, protocol.ToArray(configParams))
}

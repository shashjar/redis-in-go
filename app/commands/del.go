package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// DEL command
func del(conn net.Conn, command []string) {
	if len(command) < 2 {
		write(conn, protocol.ToSimpleError("ERR no keys for deletion provided to 'del' command"))
		return
	}

	numDeleted := 0
	for _, keyToDelete := range command[1:] {
		deleted := store.DeleteKey(keyToDelete)
		if deleted {
			numDeleted += 1
		}
	}

	write(conn, protocol.ToInteger(numDeleted))
}

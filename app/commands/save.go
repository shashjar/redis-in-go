package commands

import (
	"log"
	"net"

	"github.com/shashjar/redis-in-go/app/persistence"
	"github.com/shashjar/redis-in-go/app/protocol"
)

// SAVE command
func save(conn net.Conn) {
	err := persistence.DumpToRDB()
	if err != nil {
		log.Println("Error creating RDB file to write state of key-value store to:", err.Error())
		write(conn, protocol.ToSimpleError("ERR failed to persist state of key-value store to RDB file"))
		return
	}

	write(conn, protocol.ToSimpleString("OK"))
}

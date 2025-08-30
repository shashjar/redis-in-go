package commands

import (
	"log"
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// RPUSH command
func rpush(conn net.Conn, command []string) {
	if len(command) < 3 {
		log.Println("ERR wrong number of arguments for 'rpush' command")
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'rpush' command"))
		return
	}

	listKey := command[1]
	elements := command[2:]
	log.Println(listKey, elements)

	newListLength, errorResponse, ok := store.RPush(listKey, elements)
	if ok {
		log.Println(newListLength)
		write(conn, protocol.ToInteger(newListLength))
	} else {
		log.Println(errorResponse)
		write(conn, protocol.ToSimpleError(errorResponse))
	}
}

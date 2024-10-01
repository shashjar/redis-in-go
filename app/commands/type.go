package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// TYPE command
func typeCommand(conn net.Conn, command []string) {
	if len(command) != 2 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'type' command"))
		return
	}

	typeName, ok := store.Type(command[1])
	if ok {
		write(conn, protocol.ToSimpleString(typeName))
	} else {
		write(conn, protocol.ToSimpleString("none"))
	}
}

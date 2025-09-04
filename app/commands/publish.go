package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/pubsub"
)

// PUBLISH command
func publish(conn net.Conn, command []string) {
	if len(command) != 3 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'publish' command"))
		return
	}

	pubChan := command[1]
	pubMsg := command[2]

	numClientsPublishedTo := pubsub.Publish(pubChan, pubMsg)
	write(conn, protocol.ToInteger(numClientsPublishedTo))
}

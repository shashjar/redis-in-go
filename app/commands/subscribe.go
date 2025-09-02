package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/pubsub"
)

// SUBSCRIBE command
func subscribe(conn net.Conn, command []string) {
	if len(command) < 2 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'subscribe' command"))
		return
	}

	subChans := command[1:]
	for _, subChan := range subChans {
		numSubscriptions := pubsub.Subscribe(conn, subChan)
		write(conn, protocol.ToMixedArray([]interface{}{"subscribe", subChan, numSubscriptions}))
	}
}

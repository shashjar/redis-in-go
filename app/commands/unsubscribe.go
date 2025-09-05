package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/pubsub"
)

// UNSUBSCRIBE command
func unsubscribe(conn net.Conn, command []string) {
	if len(command) < 1 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'unsubscribe' command"))
		return
	}

	var unsubChans []string
	if len(command) == 1 {
		unsubChans = pubsub.GetSubscribedChannels(conn)
	} else {
		unsubChans = command[1:]
	}

	for _, unsubChan := range unsubChans {
		numSubscriptions := pubsub.Unsubscribe(conn, unsubChan)
		write(conn, protocol.ToMixedArray([]interface{}{"unsubscribe", unsubChan, numSubscriptions}))
	}
}

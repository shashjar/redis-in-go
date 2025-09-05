package pubsub

import (
	"log"
	"net"
	"sync"

	"github.com/shashjar/redis-in-go/app/protocol"
)

// Represents a manager of Redis pubsub functionality, mapping between clients and the
// channels they are subscribed to, and between channels and the clients that are subscribed
// to them
type PubSubManager struct {
	mu             sync.RWMutex
	clientChannels map[net.Conn]map[string]struct{} // client -> set of channels
	channelClients map[string]map[net.Conn]struct{} // channel -> set of clients
}

// Maps between clients and channels for pubsub
var PUBSUB_MANAGER = PubSubManager{
	clientChannels: make(map[net.Conn]map[string]struct{}),
	channelClients: make(map[string]map[net.Conn]struct{}),
}

// Stores clients that are in subscribed mode
var SUBSCRIBED_MODE = map[net.Conn]struct{}{}

// Stores the commands that are allowed while a client is in subscribed mode
var SUBSCRIBED_MODE_ALLOWED_COMMANDS = map[string]struct{}{
	"subscribe":   {},
	"unsubscribe": {},
	"ping":        {},
}

func InSubscribedMode(conn net.Conn) bool {
	_, ok := SUBSCRIBED_MODE[conn]
	return ok
}

// Subscribes the given client to the given channel, returning the number of channels
// the client is currently subscribed to
func Subscribe(conn net.Conn, channel string) int {
	PUBSUB_MANAGER.mu.Lock()
	defer PUBSUB_MANAGER.mu.Unlock()

	if _, ok := PUBSUB_MANAGER.clientChannels[conn]; !ok {
		PUBSUB_MANAGER.clientChannels[conn] = make(map[string]struct{})
	}

	if !isSubscribed(conn, channel) {
		PUBSUB_MANAGER.clientChannels[conn][channel] = struct{}{}
		if _, ok := PUBSUB_MANAGER.channelClients[channel]; !ok {
			PUBSUB_MANAGER.channelClients[channel] = make(map[net.Conn]struct{})
		}
		PUBSUB_MANAGER.channelClients[channel][conn] = struct{}{}
	}

	SUBSCRIBED_MODE[conn] = struct{}{}

	return len(PUBSUB_MANAGER.clientChannels[conn])
}

// Publishes the given message to the given channel, returning the number of clients that
// the message was sent to
func Publish(channel string, message string) int {
	PUBSUB_MANAGER.mu.Lock()
	defer PUBSUB_MANAGER.mu.Unlock()

	channelClients, ok := PUBSUB_MANAGER.channelClients[channel]
	if !ok {
		return 0
	}

	for channelClient := range channelClients {
		data := protocol.ToArray([]string{"message", channel, message})
		_, err := channelClient.Write([]byte(data))
		if err != nil {
			log.Println("Error writing to subscribed client:", err.Error())
			channelClient.Close()
		}
	}

	return len(channelClients)
}

func isSubscribed(conn net.Conn, channel string) bool {
	_, okClientChannel := PUBSUB_MANAGER.clientChannels[conn][channel]
	_, okChannelClient := PUBSUB_MANAGER.channelClients[channel][conn]
	if okClientChannel != okChannelClient {
		panic("clientChannels and channelClients maps are out of sync")
	}
	return okClientChannel
}

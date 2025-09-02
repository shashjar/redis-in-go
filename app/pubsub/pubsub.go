package pubsub

import (
	"net"
	"sync"
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

// Subscribes the given client to the given channel, returning the number of channels
// the client is currently subscribed to
func Subscribe(conn net.Conn, channel string) int {
	PUBSUB_MANAGER.mu.Lock()
	defer PUBSUB_MANAGER.mu.Unlock()

	if _, ok := PUBSUB_MANAGER.clientChannels[conn]; !ok {
		PUBSUB_MANAGER.clientChannels[conn] = make(map[string]struct{})
	}

	if !alreadySubscribed(conn, channel) {
		PUBSUB_MANAGER.clientChannels[conn][channel] = struct{}{}
		if _, ok := PUBSUB_MANAGER.channelClients[channel]; !ok {
			PUBSUB_MANAGER.channelClients[channel] = make(map[net.Conn]struct{})
		}
		PUBSUB_MANAGER.channelClients[channel][conn] = struct{}{}
	}

	return len(PUBSUB_MANAGER.clientChannels[conn])
}

func alreadySubscribed(conn net.Conn, channel string) bool {
	_, okClientChannel := PUBSUB_MANAGER.clientChannels[conn][channel]
	_, okChannelClient := PUBSUB_MANAGER.channelClients[channel][conn]
	if okClientChannel != okChannelClient {
		panic("clientChannels and channelClients maps are out of sync")
	}
	return okClientChannel
}

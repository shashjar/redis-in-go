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
var SUBSCRIBED_MODE_MU sync.RWMutex

// Stores the commands that are allowed while a client is in subscribed mode
var SUBSCRIBED_MODE_ALLOWED_COMMANDS = map[string]struct{}{
	"subscribe":   {},
	"unsubscribe": {},
	"ping":        {},
}

func InSubscribedMode(conn net.Conn) bool {
	SUBSCRIBED_MODE_MU.RLock()
	defer SUBSCRIBED_MODE_MU.RUnlock()

	_, ok := SUBSCRIBED_MODE[conn]
	return ok
}

// Subscribes the given client to the given channel, returning the number of channels
// the client is currently subscribed to (after subscribing)
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

	SUBSCRIBED_MODE_MU.Lock()
	SUBSCRIBED_MODE[conn] = struct{}{}
	SUBSCRIBED_MODE_MU.Unlock()

	return len(PUBSUB_MANAGER.clientChannels[conn])
}

// Unsubscribes the given client from the given channel, returning the number of channels
// the client is currently subscribed to (after unsubscribing)
func Unsubscribe(conn net.Conn, channel string) int {
	PUBSUB_MANAGER.mu.Lock()
	defer PUBSUB_MANAGER.mu.Unlock()

	if isSubscribed(conn, channel) {
		delete(PUBSUB_MANAGER.clientChannels[conn], channel)
		if len(PUBSUB_MANAGER.clientChannels[conn]) == 0 {
			delete(PUBSUB_MANAGER.clientChannels, conn)
		}

		delete(PUBSUB_MANAGER.channelClients[channel], conn)
		if len(PUBSUB_MANAGER.channelClients[channel]) == 0 {
			delete(PUBSUB_MANAGER.channelClients, channel)
		}
	}

	clientChannels, ok := PUBSUB_MANAGER.clientChannels[conn]
	if !ok {
		SUBSCRIBED_MODE_MU.Lock()
		delete(SUBSCRIBED_MODE, conn)
		SUBSCRIBED_MODE_MU.Unlock()
		return 0
	}

	return len(clientChannels)
}

// Publishes the given message to the given channel, returning the number of clients that
// the message was sent to
func Publish(channel string, message string) int {
	PUBSUB_MANAGER.mu.RLock()
	defer PUBSUB_MANAGER.mu.RUnlock()

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

// Returns all channels that the given client is currently subscribed to
func GetSubscribedChannels(conn net.Conn) []string {
	PUBSUB_MANAGER.mu.RLock()
	defer PUBSUB_MANAGER.mu.RUnlock()

	clientChannels, ok := PUBSUB_MANAGER.clientChannels[conn]
	if !ok {
		return []string{}
	}

	subbedChannels := make([]string, 0, len(clientChannels))
	for clientChannel := range clientChannels {
		subbedChannels = append(subbedChannels, clientChannel)
	}

	return subbedChannels
}

func isSubscribed(conn net.Conn, channel string) bool {
	_, okClientChannel := PUBSUB_MANAGER.clientChannels[conn][channel]
	_, okChannelClient := PUBSUB_MANAGER.channelClients[channel][conn]
	if okClientChannel != okChannelClient {
		panic("clientChannels and channelClients maps are out of sync")
	}
	return okClientChannel
}

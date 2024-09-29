package commands

import (
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// KEYS command
func keys(conn net.Conn, command []string) {
	if len(command) != 2 || command[1] != "*" {
		write(conn, protocol.ToSimpleError("ERR invalid arguments for 'keys' command"))
		return
	}

	var keys []string
	for key, value := range store.REDIS_STORE.Data {
		if value.IsExpired() {
			store.REDIS_STORE.DeleteKey(key)
		} else {
			keys = append(keys, key)
		}
	}

	write(conn, protocol.ToArray(keys))
}

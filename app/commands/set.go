package commands

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// SET command
func set(conn net.Conn, command []string) {
	if len(command) != 3 && len(command) != 5 {
		write(conn, protocol.ToSimpleError(("ERR wrong number of arguments for 'set' command")))
		return
	}

	expiresAt := time.Time{}
	if len(command) == 5 {
		var ttlUnit time.Duration
		expiryFormat := strings.ToLower(command[3])

		if expiryFormat == "ex" {
			ttlUnit = time.Second
		} else if expiryFormat == "px" {
			ttlUnit = time.Millisecond
		} else {
			write(conn, protocol.ToSimpleError(("ERR invalid expiration format for 'set' command")))
			return
		}

		ttl, err := strconv.Atoi(command[4])
		if err != nil {
			write(conn, protocol.ToSimpleError(("ERR invalid TTL value provided for 'set' command")))
			return
		}

		expiresAt = time.Now().Add(time.Duration(ttl) * ttlUnit)
	}

	store.Set(command[1], command[2], expiresAt)
	write(conn, protocol.ToSimpleString("OK"))
}

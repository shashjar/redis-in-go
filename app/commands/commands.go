package commands

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/shashjar/redis-in-go/app/persistence"
	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/replication"
	"github.com/shashjar/redis-in-go/app/store"
)

// TODO: returns an empty array to get redis-cli to initialize properly
// COMMAND DOCS command
func commandDocs(conn net.Conn) {
	write(conn, "*0\r\n")
}

// CONFIG GET command
func configGet(conn net.Conn, command []string) {
	if len(command) <= 2 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'config get' command"))
		return
	}

	var configParams []string
	for i := 2; i < len(command); i++ {
		switch command[i] {
		case "dir":
			configParams = append(configParams, "dir")
			configParams = append(configParams, persistence.RDB_DIR)
		case "dbfilename":
			configParams = append(configParams, "dbfilename")
			configParams = append(configParams, persistence.RDB_FILENAME)
		default:
			write(conn, protocol.ToSimpleError("ERR invalid configuration parameter for 'config get' command"))
			return
		}
	}

	write(conn, protocol.ToArray(configParams))
}

// INFO REPLICATION command
func infoReplication(conn net.Conn) {
	var replicationInfo string
	if store.SERVER_CONFIG.IsReplica {
		replicationInfo = "role:slave\n"
	} else {
		replicationInfo = fmt.Sprintf("role:master\nmaster_replid:%s\nmaster_repl_offset:%d\n", store.SERVER_CONFIG.MasterReplicationID, store.SERVER_CONFIG.MasterReplicationOffset)
	}

	write(conn, protocol.ToBulkString(replicationInfo))
}

// REPLCONF command
func replconf(conn net.Conn) {
	write(conn, protocol.ToSimpleString("OK"))
}

// PSYNC command
func psync(conn net.Conn) {
	response := fmt.Sprintf("FULLRESYNC %s %d", store.SERVER_CONFIG.MasterReplicationID, store.SERVER_CONFIG.MasterReplicationOffset)
	write(conn, protocol.ToSimpleString(response))
	replication.ExecuteFullResync(conn)
}

// SAVE command
func save(conn net.Conn) {
	err := persistence.DumpToRDB()
	if err != nil {
		log.Println("Error creating RDB file to write state of key-value store to:", err.Error())
		write(conn, protocol.ToSimpleError("ERR failed to persist state of key-value store to RDB file"))
		return
	}

	write(conn, protocol.ToSimpleString("OK"))
}

// PING command
func ping(conn net.Conn) {
	write(conn, protocol.ToSimpleString("PONG"))
}

// ECHO command
func echo(conn net.Conn, command []string) {
	if len(command) <= 1 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'echo' command"))
		return
	}

	write(conn, protocol.ToBulkString(command[1]))
}

// GET command
func get(conn net.Conn, command []string) {
	if len(command) != 2 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'get' command"))
		return
	}

	val, ok := store.REDIS_STORE.Get(command[1])
	if ok {
		write(conn, protocol.ToBulkString(val))
	} else {
		write(conn, protocol.ToNullBulkString())
	}
}

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

	store.REDIS_STORE.Set(command[1], command[2], expiresAt)
	write(conn, protocol.ToSimpleString("OK"))
}

// DEL command
func del(conn net.Conn, command []string) {
	if len(command) < 2 {
		write(conn, protocol.ToSimpleError("ERR no keys for deletion provided to 'del' command"))
		return
	}

	numDeleted := 0
	for _, keyToDelete := range command[1:] {
		deleted := store.REDIS_STORE.DeleteKey(keyToDelete)
		if deleted {
			numDeleted += 1
		}
	}

	write(conn, protocol.ToInteger(numDeleted))
}

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

func unknownCommand(conn net.Conn, command []string) {
	log.Printf("Unknown command: %s\n", command)
	write(conn, protocol.ToSimpleError("ERR unknown command '"+command[0]+"'"))
}

func invalidCommand(conn net.Conn, command []string) {
	log.Printf("Invalid usage of command: %s\n", command)
	write(conn, protocol.ToSimpleError("ERR invalid usage of command '"+command[0]+"'"))
}

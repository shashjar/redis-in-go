package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// TODO: returns an empty array to get redis-cli to initialize properly
// COMMAND DOCS command
func commandDocs(conn net.Conn) {
	write(conn, "*0\r\n")
}

// CONFIG GET command
func configGet(conn net.Conn, command []string) {
	if len(command) <= 2 {
		write(conn, toSimpleError("ERR wrong number of arguments for 'config get' command"))
		return
	}

	var configParams []string
	for i := 2; i < len(command); i++ {
		switch command[i] {
		case "dir":
			configParams = append(configParams, "dir")
			configParams = append(configParams, RDB_DIR)
		case "dbfilename":
			configParams = append(configParams, "dbfilename")
			configParams = append(configParams, RDB_FILENAME)
		default:
			write(conn, toSimpleError("ERR invalid configuration parameter for 'config get' command"))
			return
		}
	}

	write(conn, toArray(configParams))
}

// INFO REPLICATION command
func infoReplication(conn net.Conn) {
	var replicationInfo string
	if SERVER_CONFIG.isReplica {
		replicationInfo = "role:slave\n"
	} else {
		replicationInfo = fmt.Sprintf("role:master\nmaster_replid:%s\nmaster_repl_offset:%d\n", SERVER_CONFIG.masterReplicationID, SERVER_CONFIG.masterReplicationOffset)
	}

	write(conn, toBulkString(replicationInfo))
}

// REPLCONF command
func replconf(conn net.Conn) {
	write(conn, toSimpleString("OK"))
}

// PSYNC command
func psync(conn net.Conn) {
	response := fmt.Sprintf("FULLRESYNC %s %d", SERVER_CONFIG.masterReplicationID, SERVER_CONFIG.masterReplicationOffset)
	write(conn, toSimpleString(response))
}

// PING command
func ping(conn net.Conn) {
	write(conn, toSimpleString("PONG"))
}

// ECHO command
func echo(conn net.Conn, command []string) {
	if len(command) <= 1 {
		write(conn, toSimpleError("ERR wrong number of arguments for 'echo' command"))
		return
	}

	write(conn, toBulkString(command[1]))
}

// GET command
func get(conn net.Conn, command []string) {
	if len(command) != 2 {
		write(conn, toSimpleError("ERR wrong number of arguments for 'get' command"))
		return
	}

	val, ok := REDIS_STORE.Get(command[1])
	if ok {
		write(conn, toBulkString(val))
	} else {
		write(conn, toNullBulkString())
	}
}

// SET command
func set(conn net.Conn, command []string) {
	if len(command) != 3 && len(command) != 5 {
		write(conn, toSimpleError(("ERR wrong number of arguments for 'set' command")))
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
			write(conn, toSimpleError(("ERR invalid expiration format for 'set' command")))
			return
		}

		ttl, err := strconv.Atoi(command[4])
		if err != nil {
			write(conn, toSimpleError(("ERR invalid TTL value provided for 'set' command")))
			return
		}

		expiresAt = time.Now().Add(time.Duration(ttl) * ttlUnit)
	}

	REDIS_STORE.Set(command[1], command[2], expiresAt)
	write(conn, toSimpleString("OK"))
}

// DEL command
func del(conn net.Conn, command []string) {
	if len(command) < 2 {
		write(conn, toSimpleError("ERR no keys for deletion provided to 'del' command"))
		return
	}

	numDeleted := 0
	for _, keyToDelete := range command[1:] {
		deleted := REDIS_STORE.DeleteKey(keyToDelete)
		if deleted {
			numDeleted += 1
		}
	}

	write(conn, toInteger(numDeleted))
}

// KEYS command
func keys(conn net.Conn, command []string) {
	if len(command) != 2 || command[1] != "*" {
		write(conn, toSimpleError("ERR invalid arguments for 'keys' command"))
		return
	}

	var keys []string
	for key, value := range REDIS_STORE.data {
		if value.IsExpired() {
			REDIS_STORE.DeleteKey(key)
		} else {
			keys = append(keys, key)
		}
	}

	write(conn, toArray(keys))
}

func unknownCommand(conn net.Conn, command []string) {
	log.Printf("Unknown command: %s\n", command)
	write(conn, toSimpleError("ERR unknown command '"+command[0]+"'"))
}

func invalidCommand(conn net.Conn, command []string) {
	log.Printf("Invalid usage of command: %s\n", command)
	write(conn, toSimpleError("ERR invalid usage of command '"+command[0]+"'"))
}

package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// XREAD command
func xread(conn net.Conn, command []string) {
	if len(command) != 4 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'xread' command"))
		return
	}

	if command[1] != "streams" {
		write(conn, protocol.ToSimpleError("ERR only able to read from streams with 'xread' command"))
		return
	}

	streamKey := command[2]
	ok, startMSTime, startSeqNum, errorResponse := getEntryIDParts(command[3], true)
	if !ok {
		write(conn, protocol.ToSimpleError(errorResponse))
		return
	}

	ok, entries, errorResponse := store.XRead(streamKey, startMSTime, startSeqNum)
	if !ok {
		write(conn, protocol.ToSimpleError(errorResponse))
		return
	}

	var entriesEncoded []string
	for _, entry := range entries {
		entryIDBulkString := protocol.ToBulkString(entry.ID)
		var kvsArray []string
		for k, v := range entry.KVPairs {
			kvsArray = append(kvsArray, k)
			kvsArray = append(kvsArray, v)
		}
		kvsArrayEncoded := protocol.ToArray(kvsArray)
		entryEncoded := fmt.Sprintf("%s2\r\n%s%s", protocol.ARRAY, entryIDBulkString, kvsArrayEncoded)
		entriesEncoded = append(entriesEncoded, entryEncoded)
	}

	response := strings.Join(entriesEncoded, "")
	response = fmt.Sprintf("%s1\r\n%s2\r\n%s%s%d\r\n", protocol.ARRAY, protocol.ARRAY, protocol.ToBulkString(streamKey), protocol.ARRAY, len(entries)) + response
	write(conn, response)
}

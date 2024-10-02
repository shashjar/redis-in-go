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
	if len(command) < 4 || len(command)%2 != 0 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'xread' command"))
		return
	}

	if command[1] != "streams" {
		write(conn, protocol.ToSimpleError("ERR only able to read from streams with 'xread' command"))
		return
	}

	numStreams := (len(command) - 2) / 2
	response := fmt.Sprintf("%s%d\r\n", protocol.ARRAY, numStreams)

	for i := range numStreams {
		streamKey := command[2+i]
		ok, startMSTime, startSeqNum, errorResponse := getEntryIDParts(command[2+i+numStreams], true)
		if !ok {
			write(conn, protocol.ToSimpleError(errorResponse))
			return
		}

		ok, entries, errorResponse := store.XRead(streamKey, startMSTime, startSeqNum)
		if !ok {
			write(conn, protocol.ToSimpleError(errorResponse))
			return
		}

		streamEncoded := fmt.Sprintf("%s2\r\n%s%s%d\r\n", protocol.ARRAY, protocol.ToBulkString(streamKey), protocol.ARRAY, len(entries))
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

		streamEncoded = streamEncoded + strings.Join(entriesEncoded, "")
		response += streamEncoded
	}

	write(conn, response)
}

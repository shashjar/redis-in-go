package commands

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// XREAD command
func xread(conn net.Conn, command []string) {
	if len(command) < 4 || len(command)%2 != 0 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'xread' command"))
		return
	}

	streamsIndex, numPrefixParameters, blockingTimeMS, filterEntryNewerThanTime := 1, 2, 0, time.Time{}
	if command[1] == "block" {
		streamsIndex, numPrefixParameters = 3, 4
		blockForMS, err := strconv.Atoi(command[2])
		if err != nil || blockForMS < 0 {
			write(conn, protocol.ToSimpleError("ERR invalid blocking time provided for 'xread' command"))
			return
		}
		blockingTimeMS = blockForMS
		filterEntryNewerThanTime = time.Now()
	}

	if command[streamsIndex] != "streams" {
		write(conn, protocol.ToSimpleError("ERR only able to read from streams with 'xread' command"))
		return
	}

	time.Sleep(time.Duration(blockingTimeMS) * time.Millisecond)

	numStreams := (len(command) - numPrefixParameters) / 2
	response := fmt.Sprintf("%s%d\r\n", protocol.ARRAY, numStreams)
	numTotalEntriesReturned := 0

	for i := range numStreams {
		streamKey := command[numPrefixParameters+i]
		startMSTime, startSeqNum, errorResponse, ok := getEntryIDParts(command[numPrefixParameters+i+numStreams], true)
		if !ok {
			write(conn, protocol.ToSimpleError(errorResponse))
			return
		}

		entries, errorResponse, ok := store.XRead(streamKey, startMSTime, startSeqNum, filterEntryNewerThanTime)
		if !ok {
			write(conn, protocol.ToSimpleError(errorResponse))
			return
		}
		numTotalEntriesReturned += len(entries)

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

	if numTotalEntriesReturned == 0 {
		write(conn, protocol.ToNullBulkString())
		return
	}

	write(conn, response)
}

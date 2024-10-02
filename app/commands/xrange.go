package commands

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

// XRANGE command
func xrange(conn net.Conn, command []string) {
	if len(command) != 4 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for 'xrange' command"))
		return
	}

	streamKey := command[1]
	ok, startMillisecondsTime, startSequenceNumber, errorResponse := getEntryIDParts(command[2], true)
	if !ok {
		write(conn, protocol.ToSimpleError(errorResponse))
		return
	}
	ok, endMillisecondsTime, endSequenceNumber, errorResponse := getEntryIDParts(command[3], false)
	if !ok {
		write(conn, protocol.ToSimpleError(errorResponse))
		return
	}

	ok, entries, errorResponse := store.XRange(streamKey, startMillisecondsTime, startSequenceNumber, endMillisecondsTime, endSequenceNumber)
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
	response = fmt.Sprintf("%s%d\r\n", protocol.ARRAY, len(entries)) + response
	write(conn, response)
}

func getEntryIDParts(entryID string, isStart bool) (bool, int, int, string) {
	if isStart && entryID == "-" {
		return true, 0, 0, ""
	}

	if !isStart && entryID == "+" {
		return true, math.MaxInt, math.MaxInt, ""
	}

	parts := strings.Split(entryID, "-")
	if len(parts) == 1 {
		millisecondsTime, err := strconv.Atoi(parts[0])
		if err != nil {
			return false, 0, 0, "ERR invalid millisecondsTime parameter for XRANGE command"
		}
		if isStart {
			return true, millisecondsTime, 0, ""
		} else {
			return true, millisecondsTime, math.MaxInt, ""
		}
	} else if len(parts) == 2 {
		millisecondsTime, err := strconv.Atoi(parts[0])
		if err != nil {
			return false, 0, 0, "ERR invalid millisecondsTime parameter for XRANGE command"
		}
		sequenceNumber, err := strconv.Atoi(parts[1])
		if err != nil {
			return false, 0, 0, "ERR invalid sequenceNumber parameter for XRANGE command"
		}
		return true, millisecondsTime, sequenceNumber, ""
	} else {
		return false, 0, 0, "ERR invalid ID provided to XRANGE command"
	}
}

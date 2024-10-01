package store

import (
	"strconv"
	"strings"
)

type Stream struct {
	Entries []StreamEntry
}

type StreamEntry struct {
	ID      string
	KVPairs map[string]string
}

func (stream *Stream) addEntry(entryID string, keys []string, values []string) {
	kvPairs := make(map[string]string)
	for i, key := range keys {
		kvPairs[key] = values[i]
	}

	newEntry := StreamEntry{ID: entryID, KVPairs: kvPairs}
	stream.Entries = append(stream.Entries, newEntry)
}

func (stream *Stream) validEntryID(entryID string) (bool, string) {
	ok, millisecondsTime, sequenceNumber, errorResponse := splitEntryID(entryID)
	if !ok {
		return false, errorResponse
	}

	if len(stream.Entries) == 0 {
		if (millisecondsTime >= 0 && sequenceNumber >= 0) && (millisecondsTime > 0 || sequenceNumber > 0) {
			return true, ""
		} else {
			return false, "ERR The ID specified in XADD must be greater than 0-0"
		}
	} else {
		lastEntry := stream.Entries[len(stream.Entries)-1]
		_, prevMillisecondsTime, prevSequenceNumber, _ := splitEntryID(lastEntry.ID)

		if millisecondsTime < prevMillisecondsTime {
			return false, "ERR The ID specified in XADD must have a greater millisecondsTime than the previous entry in the stream"
		}

		if (millisecondsTime == prevMillisecondsTime) && (sequenceNumber <= prevSequenceNumber) {
			return false, "ERR The ID specified in XADD must have a greater sequenceNumber than the previous entry in the stream if times are equal"
		}
	}

	return true, ""
}

func splitEntryID(entryID string) (bool, int, int, string) {
	parts := strings.Split(entryID, "-")
	if len(parts) != 2 {
		return false, 0, 0, "ERR The ID specified in XADD does not follow the correct hyphenated format"
	}

	millisecondsTime, err := strconv.Atoi(parts[0])
	if err != nil {
		return false, 0, 0, "ERR The ID specified in XADD has an invalid millisecondsTime parameter"
	}

	sequenceNumber, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, 0, 0, "ERR The ID specified in XADD has an invalid sequenceNumber parameter"
	}

	return true, millisecondsTime, sequenceNumber, ""
}

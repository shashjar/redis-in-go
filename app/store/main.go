package store

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// NOTE: this key-value store currently supports only passive expiration (an expired key is deleted only after a client attempts to access it).
// May add active expiration via period random sampling/testing in the future (https://redis.io/docs/latest/commands/expire/#how-redis-expires-keys).

// Maps string keys to the values they are associated with in the Redis store
var REDIS_STORE = KeyValueStore{data: make(map[string]KeyValue)}

// Retrieves the complete mapping of data currently in the key-value store
func Data() map[string]KeyValue {
	return REDIS_STORE.data
}

// Returns the type for the given key, if that key exists
func Type(key string) (string, bool) {
	kv, ok := REDIS_STORE.get(key)
	if !ok {
		return "", false
	}

	return kv.Type, true
}

// Gets the string value associated with the given key, if that key exists
func Get(key string) (string, bool) {
	kv, ok := REDIS_STORE.get(key)
	if !ok || kv.Type != "string" {
		return "", false
	}

	return kv.Value.(string), true
}

// Sets the string value associated with the given key, with some expiration
func Set(key string, value string, expiration time.Time) {
	REDIS_STORE.setString(key, value, expiration)
}

// Deletes the key-value pair associated with the given key
func DeleteKey(key string) bool {
	return REDIS_STORE.deleteKey(key)
}

// Increments the value associated with the given key, if it exists and is an integer value
func Incr(key string) (int, bool) {
	kv, ok := REDIS_STORE.get(key)
	if ok && kv.Type != "string" {
		return 0, false
	}

	stringInt := ""
	if ok {
		stringInt = kv.Value.(string)
	} else {
		stringInt = "0"
	}

	num, err := strconv.Atoi(stringInt)
	if err != nil {
		return 0, false
	}

	Set(key, strconv.Itoa(num+1), time.Time{})
	return num + 1, true
}

// Creates an entry in the given stream
func XAdd(streamKey string, entryID string, keys []string, values []string) (string, string, bool) {
	kv, ok := REDIS_STORE.get(streamKey)
	if !ok || kv.Type != "stream" {
		return createStream(streamKey, entryID, keys, values, &REDIS_STORE)
	} else {
		stream := kv.Value.(*Stream)
		return addEntryToStream(stream, entryID, keys, values)
	}
}

// Retrieves a range of entries in the given stream
func XRange(streamKey string, startMSTime int, startSeqNum int, endMSTime int, endSeqNum int) ([]StreamEntry, string, bool) {
	kv, ok := REDIS_STORE.get(streamKey)
	if !ok || kv.Type != "stream" {
		return nil, "ERR stream with key provided to XRANGE command not found", false
	}

	stream := kv.Value.(*Stream)
	return stream.getEntriesInRange(startMSTime, startSeqNum, endMSTime, endSeqNum, time.Time{}, false), "", true
}

// Reads entries later than some start index in the given stream
func XRead(streamKey string, startMSTime int, startSeqNum int, filterEntryNewerThanTime time.Time) ([]StreamEntry, string, bool) {
	kv, ok := REDIS_STORE.get(streamKey)
	if !ok || kv.Type != "stream" {
		return nil, "ERR stream with key provided to XREAD command not found", false
	}

	stream := kv.Value.(*Stream)
	return stream.getEntriesInRange(startMSTime, startSeqNum, math.MaxInt, math.MaxInt, filterEntryNewerThanTime, true), "", true
}

// Appends the given elements to the end of the list associated with the given key,
// creating that list if it does not exist
func RPush(listKey string, elements []string) (int, string, bool) {
	kv, ok := REDIS_STORE.get(listKey)
	if ok && kv.Type != "list" {
		return 0, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	}

	var list *List
	if !ok {
		list = createEmptyList(listKey, &REDIS_STORE)
	} else {
		list = kv.Value.(*List)
	}

	newListLength := list.appendElements(elements)
	notifyBlpopWaiter(listKey)
	return newListLength, "", true
}

// Inserts the given elements at the front d of the list associated with the given key,
// creating that list if it does not exist
func LPush(listKey string, elements []string) (int, string, bool) {
	kv, ok := REDIS_STORE.get(listKey)
	if ok && kv.Type != "list" {
		return 0, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	}

	var list *List
	if !ok {
		list = createEmptyList(listKey, &REDIS_STORE)
	} else {
		list = kv.Value.(*List)
	}

	newListLength := list.prependElements(reverseSlice(elements))
	notifyBlpopWaiter(listKey)
	return newListLength, "", true
}

// Returns the elements in the list associated with the given key, in the specified range
func LRange(listKey string, startIndex int, stopIndex int) ([]string, string, bool) {
	kv, ok := REDIS_STORE.get(listKey)
	if ok && kv.Type != "list" {
		return []string{}, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	} else if !ok {
		return []string{}, "", true
	}

	list := kv.Value.(*List)
	return list.getElementsInRange(startIndex, stopIndex), "", true
}

// Returns the length of the list associated with the given key, considering that list
// as empty if the key does not exist
func LLen(listKey string) (int, string, bool) {
	kv, ok := REDIS_STORE.get(listKey)
	if ok && kv.Type != "list" {
		return 0, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	} else if !ok {
		return 0, "", true
	}

	list := kv.Value.(*List)
	return len(list.Entries), "", true
}

// Removes and returns the first elements from the list associated with the given key,
// considering that list as empty if the key does not exist
func LPop(listKey string, popCount int) ([]string, string, bool) {
	kv, ok := REDIS_STORE.get(listKey)
	if ok && kv.Type != "list" {
		return []string{}, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	} else if !ok {
		return []string{}, "", true
	}

	list := kv.Value.(*List)
	poppedElements := list.popLeftElements(popCount)
	return poppedElements, "", true
}

// Removes and returns the first element from the first list among the given list keys that is not empty,
// blocking for up to the specified timeout if necessary
func BLPop(listKeys []string, timeoutSec int) (string, string, bool, string, bool) {
	for _, listKey := range listKeys {
		kv, ok := REDIS_STORE.get(listKey)
		if ok && kv.Type != "list" {
			return "", "", false, "WRONGTYPE Operation against a key holding the wrong kind of value", false
		} else if ok {
			list := kv.Value.(*List)
			if len(list.Entries) > 0 {
				poppedElements := list.popLeftElements(1)
				return poppedElements[0], listKey, true, "", true
			}
		}
	}

	waitChan := make(chan string, 1)
	registerBlpopWaiter(listKeys, waitChan)
	defer cleanUpBlpopWaiters(listKeys, waitChan)

	var timeoutChan <-chan time.Time
	if timeoutSec > 0 {
		timeoutChan = time.After(time.Duration(timeoutSec) * time.Second)
	}

	select {
	case listKey := <-waitChan:
		poppedElements, errorResponse, ok := LPop(listKey, 1)
		if !ok {
			return "", "", false, errorResponse, false
		}
		if len(poppedElements) != 1 {
			panic(fmt.Sprintf("Expected 1 popped element, got %d", len(poppedElements)))
		}
		return poppedElements[0], listKey, true, "", true
	case <-timeoutChan:
		return "", "", false, "", true
	}
}

// Adds all the specified members with the specified scores to the sorted set associated with the given key,
// creating it if it does not already exist
func ZAdd(setKey string, memberScores map[string]float64) (int, string, bool) {
	kv, ok := REDIS_STORE.get(setKey)
	if ok && kv.Type != "zset" {
		return 0, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	}

	var sortedSet *SortedSet
	if !ok {
		sortedSet = createEmptySortedSet(setKey, &REDIS_STORE)
	} else {
		sortedSet = kv.Value.(*SortedSet)
	}

	numNewMembers := sortedSet.addMembers(memberScores)
	return numNewMembers, "", true
}

// Returns the rank and score of the member in the sorted set associated with the given key
func ZRank(setKey string, member string) (int, float64, bool, string, bool) {
	kv, ok := REDIS_STORE.get(setKey)
	if ok && kv.Type != "zset" {
		return 0, 0, false, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	} else if !ok {
		return 0, 0, false, "", true
	}

	sortedSet := kv.Value.(*SortedSet)

	rank, score, memberExists := sortedSet.getRankAndScore(member)
	return rank, score, memberExists, "", true
}

// Returns the sorted elements in the range [startIndex, stopIndex] in the sorted set associated
// with the given key
func ZRange(setKey string, startIndex int, stopIndex int) ([]string, string, bool) {
	kv, ok := REDIS_STORE.get(setKey)
	if ok && kv.Type != "zset" {
		return nil, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	} else if !ok {
		return []string{}, "", true
	}

	sortedSet := kv.Value.(*SortedSet)
	return sortedSet.SkipList.GetElementsInRange(startIndex, stopIndex), "", true
}

// Returns the cardinality (number of elements) of the sorted set associated with the given key
func ZCard(setKey string) (int, string, bool) {
	kv, ok := REDIS_STORE.get(setKey)
	if ok && kv.Type != "zset" {
		return 0, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	} else if !ok {
		return 0, "", true
	}

	sortedSet := kv.Value.(*SortedSet)
	return sortedSet.SkipList.size, "", true
}

// Returns the score of the member in the sorted set associated with the given key
func ZScore(setKey string, member string) (float64, bool, string, bool) {
	kv, ok := REDIS_STORE.get(setKey)
	if ok && kv.Type != "zset" {
		return 0, false, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	} else if !ok {
		return 0, false, "", true
	}

	sortedSet := kv.Value.(*SortedSet)

	score, memberExists := sortedSet.Scores[member]
	return score, memberExists, "", true
}

// Removes the specified members from the sorted set associated with the given key,
// returning the number of members removed
func ZRem(setKey string, members []string) (int, string, bool) {
	kv, ok := REDIS_STORE.get(setKey)
	if ok && kv.Type != "zset" {
		return 0, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	} else if !ok {
		return 0, "", true
	}

	sortedSet := kv.Value.(*SortedSet)
	return sortedSet.removeMembers(members), "", true
}

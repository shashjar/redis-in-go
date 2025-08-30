package store

import (
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
	createdEntryID, errorResponse, ok := REDIS_STORE.xadd(streamKey, entryID, keys, values)
	if !ok {
		return "", errorResponse, false
	}

	return createdEntryID, "", true
}

// Retrieves a range of entries in the given stream
func XRange(streamKey string, startMSTime int, startSeqNum int, endMSTime int, endSeqNum int) ([]StreamEntry, string, bool) {
	entries, errorResponse, ok := REDIS_STORE.xrange(streamKey, startMSTime, startSeqNum, endMSTime, endSeqNum)
	if !ok {
		return nil, errorResponse, false
	}

	return entries, "", true
}

// Reads entries later than some start index in the given stream
func XRead(streamKey string, startMSTime int, startSeqNum int, filterEntryNewerThanTime time.Time) ([]StreamEntry, string, bool) {
	entries, errorResponse, ok := REDIS_STORE.xread(streamKey, startMSTime, startSeqNum, filterEntryNewerThanTime)
	if !ok {
		return nil, errorResponse, false
	}

	return entries, "", true
}

// Appends the given elements to the end of the list associated with the given key,
// creating that list if it does not exist
func RPush(listKey string, elements []string) (int, string, bool) {
	return REDIS_STORE.rpush(listKey, elements)
}

// Inserts the given elements at the front d of the list associated with the given key,
// creating that list if it does not exist
func LPush(listKey string, elements []string) (int, string, bool) {
	return REDIS_STORE.lpush(listKey, elements)
}

// Returns the elements in the list associated with the given key, in the specified range
func LRange(listKey string, startIndex int, stopIndex int) ([]string, string, bool) {
	return REDIS_STORE.lrange(listKey, startIndex, stopIndex)
}

// Returns the length of the list associated with the given key, considering that list
// as empty if the key does not exist
func LLen(listKey string) (int, string, bool) {
	return REDIS_STORE.llen(listKey)
}

// Removes and returns the first elements from the list associated with the given key,
// considering that list as empty if the key does not exist
func LPop(listKey string, popCount int) ([]string, string, bool) {
	return REDIS_STORE.lpop(listKey, popCount)
}

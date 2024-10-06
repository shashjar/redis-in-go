package store

import (
	"strconv"
	"time"
)

// NOTE: this key-value store currently supports only passive expiration (an expired key is deleted only after a client attempts to access it).
// May add active expiration via period random sampling/testing in the future (https://redis.io/docs/latest/commands/expire/#how-redis-expires-keys).

var REDIS_STORE = KeyValueStore{data: make(map[string]KeyValue)}

func Data() map[string]KeyValue {
	return REDIS_STORE.data
}

func Type(key string) (string, bool) {
	kv, ok := REDIS_STORE.get(key)
	if !ok {
		return "", false
	}

	return kv.Type, true
}

func Get(key string) (string, bool) {
	kv, ok := REDIS_STORE.get(key)
	if !ok || kv.Type != "string" {
		return "", false
	}

	return kv.Value.(string), true
}

func Set(key string, value string, expiration time.Time) {
	REDIS_STORE.setString(key, value, expiration)
}

func DeleteKey(key string) bool {
	return REDIS_STORE.deleteKey(key)
}

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

func XAdd(streamKey string, entryID string, keys []string, values []string) (string, string, bool) {
	createdEntryID, errorResponse, ok := REDIS_STORE.xadd(streamKey, entryID, keys, values)
	if !ok {
		return "", errorResponse, false
	}

	return createdEntryID, "", true
}

func XRange(streamKey string, startMSTime int, startSeqNum int, endMSTime int, endSeqNum int) ([]StreamEntry, string, bool) {
	entries, errorResponse, ok := REDIS_STORE.xrange(streamKey, startMSTime, startSeqNum, endMSTime, endSeqNum)
	if !ok {
		return nil, errorResponse, false
	}

	return entries, "", true
}

func XRead(streamKey string, startMSTime int, startSeqNum int, filterEntryNewerThanTime time.Time) ([]StreamEntry, string, bool) {
	entries, errorResponse, ok := REDIS_STORE.xread(streamKey, startMSTime, startSeqNum, filterEntryNewerThanTime)
	if !ok {
		return nil, errorResponse, false
	}

	return entries, "", true
}

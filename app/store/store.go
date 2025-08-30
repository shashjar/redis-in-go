package store

import (
	"math"
	"sync"
	"time"
)

// Represents a store of key-value pairs stored as a mapping, along with a mutex for read/write locking
type KeyValueStore struct {
	data map[string]KeyValue
	mu   sync.RWMutex
}

func (kvs *KeyValueStore) get(key string) (KeyValue, bool) {
	kvs.mu.RLock()
	defer kvs.mu.RUnlock()

	kv, ok := kvs.data[key]
	if !ok {
		return kv, false
	}

	if kv.IsExpired() {
		kvs.deleteKey(key)
		return kv, false
	}

	return kv, true
}

func (kvs *KeyValueStore) setString(key string, value string, expiration time.Time) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	kvs.data[key] = KeyValue{Value: value, Type: "string", Expiration: expiration}
}

// Deletes the provided key from the store. Is a no-op if the key does not exist in the store.
// Returns a boolean indicating whether the key existed and was deleted.
func (kvs *KeyValueStore) deleteKey(key string) bool {
	_, ok := kvs.data[key]
	delete(kvs.data, key)
	return ok
}

func (kvs *KeyValueStore) xadd(streamKey string, entryID string, keys []string, values []string) (string, string, bool) {
	kv, ok := kvs.get(streamKey)
	if !ok || kv.Type != "stream" {
		return createStream(streamKey, entryID, keys, values, kvs)
	} else {
		stream := kv.Value.(*Stream)
		return addEntryToStream(stream, entryID, keys, values)
	}
}

func (kvs *KeyValueStore) xrange(streamKey string, startMSTime int, startSeqNum int, endMSTime int, endSeqNum int) ([]StreamEntry, string, bool) {
	kv, ok := kvs.get(streamKey)
	if !ok || kv.Type != "stream" {
		return nil, "ERR stream with key provided to XRANGE command not found", false
	}

	stream := kv.Value.(*Stream)
	return stream.getEntriesInRange(startMSTime, startSeqNum, endMSTime, endSeqNum, time.Time{}, false), "", true
}

func (kvs *KeyValueStore) xread(streamKey string, startMSTime int, startSeqNum int, filterEntryNewerThanTime time.Time) ([]StreamEntry, string, bool) {
	kv, ok := kvs.get(streamKey)
	if !ok || kv.Type != "stream" {
		return nil, "ERR stream with key provided to XREAD command not found", false
	}

	stream := kv.Value.(*Stream)
	return stream.getEntriesInRange(startMSTime, startSeqNum, math.MaxInt, math.MaxInt, filterEntryNewerThanTime, true), "", true
}

func (kvs *KeyValueStore) rpush(listKey string, elements []string) (int, string, bool) {
	kv, ok := kvs.get(listKey)
	if ok && kv.Type != "list" {
		return 0, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	}

	var list *List
	if !ok {
		list = createEmptyList(listKey, kvs)
	} else {
		list = kv.Value.(*List)
	}

	newListLength := list.appendElements(elements)
	return newListLength, "", true
}

func (kvs *KeyValueStore) lrange(listKey string, startIndex int, stopIndex int) ([]string, string, bool) {
	kv, ok := kvs.get(listKey)
	if ok && kv.Type != "list" {
		return []string{}, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	} else if !ok {
		return []string{}, "", true
	}

	list := kv.Value.(*List)
	return list.getElementsInRange(startIndex, stopIndex), "", true
}

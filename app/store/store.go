package store

import (
	"sync"
	"time"
)

type KeyValueStore struct {
	data map[string]KeyValue
	mu   sync.RWMutex
}

// TODO: currently this only does passive expiration (an expired key is deleted only
// after a client attempts to access it). Implement active expiration as a challenge:
// https://redis.io/docs/latest/commands/expire/#how-redis-expires-keys
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

func (kvs *KeyValueStore) xadd(streamKey string, entryID string, keys []string, values []string) (bool, string, string) {
	kv, ok := kvs.get(streamKey)
	if !ok || kv.Type != "stream" {
		return createStream(streamKey, entryID, keys, values, kvs)
	} else {
		stream := kv.Value.(*Stream)
		return addEntryToStream(stream, entryID, keys, values)
	}
}

func (kvs *KeyValueStore) xrange(streamKey string, startMillisecondsTime int, startSequenceNumber int, endMillisecondsTime int, endSequenceNumber int) (bool, []StreamEntry, string) {
	kv, ok := kvs.get(streamKey)
	if !ok || kv.Type != "stream" {
		return false, nil, "ERR stream with key provided to XRANGE command not found"
	}

	stream := kv.Value.(*Stream)
	return true, stream.getEntriesInRange(startMillisecondsTime, startSequenceNumber, endMillisecondsTime, endSequenceNumber), ""
}

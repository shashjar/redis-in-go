package store

import (
	"sync"
	"time"
)

var REDIS_STORE = KeyValueStore{Data: make(map[string]KeyValue)}

type KeyValue struct {
	Value      string
	Expiration time.Time
}

type KeyValueStore struct {
	Data map[string]KeyValue
	mu   sync.RWMutex
}

func (kv *KeyValue) HasExpiration() bool {
	return !kv.Expiration.IsZero()
}

func (kv *KeyValue) IsExpired() bool {
	return kv.HasExpiration() && time.Now().After(kv.Expiration)
}

// TODO: currently this only does passive expiration (an expired key is deleted only
// after a client attempts to access it). Implement active expiration as a challenge:
// https://redis.io/docs/latest/commands/expire/#how-redis-expires-keys
func (kvs *KeyValueStore) Get(key string) (string, bool) {
	kvs.mu.RLock()
	defer kvs.mu.RUnlock()

	item, ok := kvs.Data[key]

	if !ok {
		return "", false
	}

	if item.IsExpired() {
		kvs.DeleteKey(key)
		return "", false
	}

	return item.Value, true
}

func (kvs *KeyValueStore) Set(key string, value string, expiration time.Time) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	kvs.Data[key] = KeyValue{Value: value, Expiration: expiration}
}

// Deletes the provided key from the store. Is a no-op if the key does not exist in the store.
// Returns a boolean indicating whether the key existed and was deleted.
func (kvs *KeyValueStore) DeleteKey(key string) bool {
	_, ok := kvs.Data[key]
	delete(kvs.Data, key)
	return ok
}

func (kvs *KeyValueStore) GetData() map[string]KeyValue {
	return kvs.Data
}

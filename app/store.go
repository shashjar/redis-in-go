package main

import (
	"sync"
	"time"
)

var REDIS_STORE = KeyValueStore{data: make(map[string]KeyValue)}

type KeyValue struct {
	value      string
	expiration time.Time
}

type KeyValueStore struct {
	data map[string]KeyValue
	mu   sync.RWMutex
}

func (kv *KeyValue) IsExpired() bool {
	return !kv.expiration.IsZero() && time.Now().After(kv.expiration)
}

// TODO: currently this only does passive expiration (an expired key is deleted only
// after a client attempts to access it). Implement active expiration as a challenge:
// https://redis.io/docs/latest/commands/expire/#how-redis-expires-keys
func (kvs *KeyValueStore) Get(key string) (string, bool) {
	kvs.mu.RLock()
	defer kvs.mu.RUnlock()

	item, ok := kvs.data[key]

	if !ok {
		return "", false
	}

	if item.IsExpired() {
		kvs.DeleteKey(key)
		return "", false
	}

	return item.value, true
}

func (kvs *KeyValueStore) Set(key string, value string, expiration time.Time) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	kvs.data[key] = KeyValue{value: value, expiration: expiration}
}

// Deletes the provided key from the store. Is a no-op if the key does not exist in the store.
// Returns a boolean indicating whether the key existed and was deleted.
func (kvs *KeyValueStore) DeleteKey(key string) bool {
	_, ok := kvs.data[key]
	delete(kvs.data, key)
	return ok
}

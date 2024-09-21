package main

import (
	"sync"
	"time"
)

type KeyValue struct {
	value      string
	expiration time.Time
}

type KeyValueStore struct {
	data map[string]KeyValue
	mu   sync.RWMutex
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

	if item.expiration.IsZero() || time.Now().Before(item.expiration) {
		return item.value, true
	}

	delete(kvs.data, key)
	return "", false
}

func (kvs *KeyValueStore) Set(key string, value string, expiration time.Time) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	kvs.data[key] = KeyValue{value: value, expiration: expiration}
}

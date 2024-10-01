package store

import (
	"time"
)

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

func XAdd(streamKey string, entryID string, keys []string, values []string) (bool, string) {
	ok, errorResponse := REDIS_STORE.xadd(streamKey, entryID, keys, values)
	if !ok {
		return false, errorResponse
	}

	return true, ""
}

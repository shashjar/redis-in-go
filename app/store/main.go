package store

import "time"

var REDIS_STORE = KeyValueStore{data: make(map[string]KeyValue)}

func Get(key string) (string, bool) {
	return REDIS_STORE.get(key)
}

func Set(key string, value string, expiration time.Time) {
	REDIS_STORE.set(key, value, expiration)
}

func DeleteKey(key string) bool {
	return REDIS_STORE.deleteKey(key)
}

func Data() map[string]KeyValue {
	return REDIS_STORE.data
}

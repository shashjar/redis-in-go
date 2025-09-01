package store

import (
	"fmt"
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
	notifyBlpopWaiter(listKey)
	return newListLength, "", true
}

func (kvs *KeyValueStore) lpush(listKey string, elements []string) (int, string, bool) {
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

	newListLength := list.prependElements(reverseSlice(elements))
	notifyBlpopWaiter(listKey)
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

func (kvs *KeyValueStore) llen(listKey string) (int, string, bool) {
	kv, ok := kvs.get(listKey)
	if ok && kv.Type != "list" {
		return 0, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	} else if !ok {
		return 0, "", true
	}

	list := kv.Value.(*List)
	return len(list.Entries), "", true
}

func (kvs *KeyValueStore) lpop(listKey string, popCount int) ([]string, string, bool) {
	kv, ok := kvs.get(listKey)
	if ok && kv.Type != "list" {
		return []string{}, "WRONGTYPE Operation against a key holding the wrong kind of value", false
	} else if !ok {
		return []string{}, "", true
	}

	list := kv.Value.(*List)
	poppedElements := list.popLeftElements(popCount)
	return poppedElements, "", true
}

func (kvs *KeyValueStore) blpop(listKeys []string, timeoutSec int) (string, string, bool, string, bool) {
	for _, listKey := range listKeys {
		kv, ok := kvs.get(listKey)
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
		poppedElements, errorResponse, ok := kvs.lpop(listKey, 1)
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

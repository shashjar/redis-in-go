package store

import "time"

// Represents a key-value store value associated with some key
type KeyValue struct {
	Value      interface{}
	Type       string
	Expiration time.Time
}

// Determines whether the given KV has an expiration
func (kv *KeyValue) HasExpiration() bool {
	return !kv.Expiration.IsZero()
}

// Determines whether the given KV is already expired
func (kv *KeyValue) IsExpired() bool {
	return kv.HasExpiration() && time.Now().After(kv.Expiration)
}

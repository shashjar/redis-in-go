package store

import "time"

type KeyValue struct {
	Value      string
	Expiration time.Time
}

func (kv *KeyValue) HasExpiration() bool {
	return !kv.Expiration.IsZero()
}

func (kv *KeyValue) IsExpired() bool {
	return kv.HasExpiration() && time.Now().After(kv.Expiration)
}

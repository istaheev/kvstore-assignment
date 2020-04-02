package kvstore

import "sync"

// KeyValueStore is a place where key/value pairs can be stored
type KeyValueStore struct {
	// There are more readers than writers, let's try to not block readers
	// without the reason
	mu    sync.RWMutex
	kvmap map[string]string
}

// Get retrieves a key from the store.
// It returns the key's value (empty if not found) and a boolean value telling
// if the key was found in the store
func (kvs *KeyValueStore) Get(key string) (string, bool) {
	kvs.mu.RLock()
	defer kvs.mu.RUnlock()
	value, ok := kvs.kvmap[key]
	return value, ok
}

// Set sets a key and its value in the store.
func (kvs *KeyValueStore) Set(key, value string) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()
	kvs.kvmap[key] = value
}

// New creates new empty key/value store
func New() KeyValueStore {
	return KeyValueStore{
		kvmap: make(map[string]string),
	}
}

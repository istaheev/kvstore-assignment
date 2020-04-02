package kvstore

// KeyValueStore is a place where key/value pairs can be stored
type KeyValueStore struct {
	kvmap map[string]string
}

// Get retrieves a key from the store.
// It returns the key's value (empty if not found) and a boolean value telling
// if the key was found in the store
func (kvs *KeyValueStore) Get(key string) (string, bool) {
	value, ok := kvs.kvmap[key]
	return value, ok
}

// Set sets a key and its value in the store.
func (kvs *KeyValueStore) Set(key, value string) {
	kvs.kvmap[key] = value
}

// New creates new empty key/value store
func New() KeyValueStore {
	return KeyValueStore{
		kvmap: make(map[string]string),
	}
}

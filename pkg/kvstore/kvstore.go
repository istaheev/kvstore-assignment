package kvstore

// KeyValueStore is a place where key/value pairs can be stored
type KeyValueStore struct {
}

// New creates new empty key/value store
func New() KeyValueStore {
	return KeyValueStore{}
}

package kvstore

import (
	"container/list"
	"sync"
	"time"
)

type item struct {
	key       string
	value     string
	expiresAt time.Time
}

// KeyValueStore is a place where key/value pairs can be stored.
// Two data structures are used to store elements:
// - Double-linked list with all the items in the age order (from newest to
//   oldest). Using of double-linked list allows us to avoid excessive memory
//   copying when we add or remove new elements thus making it O(1) instead of
//   O(n);
// - A key map mapping keys to corresponding elements of the list above. Its
//   purpose to speed up lookup of keys;
type KeyValueStore struct {
	// There are more readers than writers, let's try to not block readers
	// without the reason
	mu          sync.RWMutex
	expireAfter time.Duration
	keysIndex   map[string]*list.Element
	items       *list.List
}

// SetExpiration sets time after which keys are expired.
// Zero value means no expiration.
func (kvs *KeyValueStore) SetExpiration(expireAfter time.Duration) {
	kvs.expireAfter = expireAfter
}

// Get retrieves a key from the store.
// It returns the key's value (empty if not found) and a boolean value telling
// if the key was found in the store.
// The complexity of the function should be no more than complexity of the
// map access, most likely amortized O(1).
func (kvs *KeyValueStore) Get(key string) (string, bool) {
	kvs.mu.RLock()
	defer kvs.mu.RUnlock()

	elem, ok := kvs.keysIndex[key]
	if !ok {
		return "", false
	}

	it := elem.Value.(item)

	// If the key is expired then make it unavailable for the caller
	if kvs.expireAfter > 0 && it.expiresAt.Before(time.Now()) {
		return "", false
	}

	return it.value, true
}

// Set sets a key and its value in the store.
// The complexity of the function should be no more than complexity of the
// map access, most likely amortized O(1).
func (kvs *KeyValueStore) Set(key, value string) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()
	expiresAt := time.Now().Add(kvs.expireAfter)

	// If there is an existing item remove it
	existing, ok := kvs.keysIndex[key]
	if ok {
		kvs.items.Remove(existing)
	}

	elem := kvs.items.PushFront(item{key, value, expiresAt})
	kvs.keysIndex[key] = elem
}

// RemoveExpired removes expired items from the store.
// No more than maxItems are removed in a single operation. Setting
// maxItems to some reasonable, not so big, number allows to keep
// the complexity of this operation close to O(1) while making it fast
// enough.
// It returns amount of keys removed.
func (kvs *KeyValueStore) RemoveExpired(maxItems int) int {
	if kvs.expireAfter == 0 {
		return 0
	}

	kvs.mu.Lock()
	defer kvs.mu.Unlock()

	now := time.Now()
	removed := 0

	for ; kvs.items.Len() > 0 && maxItems > 0; maxItems-- {
		lastElem := kvs.items.Back()
		it := lastElem.Value.(item)

		if it.expiresAt.Before(now) {
			// This item is expired
			delete(kvs.keysIndex, it.key)
			kvs.items.Remove(lastElem)
			removed++
		}
	}

	return removed
}

// New creates new empty key/value store
func New() KeyValueStore {
	return KeyValueStore{
		keysIndex: make(map[string]*list.Element),
		items:     list.New(),
	}
}

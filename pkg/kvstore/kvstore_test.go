package kvstore

import (
	"testing"
	"time"
)

func TestGetSet(t *testing.T) {
	kvs := New()

	kvs.Set("key1", "first")
	kvs.Set("key2", "second")

	shouldHave(t, &kvs, "key1", "first")
	shouldHave(t, &kvs, "key2", "second")
	shouldNotHave(t, &kvs, "key3")
}

func TestOverwrite(t *testing.T) {
	kvs := New()

	kvs.Set("key1", "first")
	kvs.Set("key1", "another first")

	shouldHave(t, &kvs, "key1", "another first")
}

func TestExpiration(t *testing.T) {
	kvs := New()
	kvs.SetExpiration(50 * time.Millisecond)

	kvs.Set("key1", "first")
	kvs.Set("key2", "second")

	time.Sleep(25 * time.Millisecond)

	kvs.Set("key1", "another first")
	kvs.Set("key3", "third")

	time.Sleep(30 * time.Millisecond)

	shouldHave(t, &kvs, "key1", "another first")
	shouldHave(t, &kvs, "key3", "third")
	shouldNotHave(t, &kvs, "key2")

	expectedRemoved := 1
	removed := kvs.RemoveExpired(10)
	if removed != expectedRemoved {
		t.Errorf("RemoveExpired: expected %d, got %d", expectedRemoved, removed)
	}

	shouldHave(t, &kvs, "key1", "another first")
	shouldHave(t, &kvs, "key3", "third")
	shouldNotHave(t, &kvs, "key2")
}

func shouldHave(t *testing.T, kvs *KeyValueStore, key, expectedValue string) {
	value, ok := kvs.Get(key)
	if !ok {
		t.Errorf("Key '%s' not found", key)
	}

	if value != expectedValue {
		t.Errorf("Key '%s' has value '%s', expected '%s'", key, value, expectedValue)
	}
}

func shouldNotHave(t *testing.T, kvs *KeyValueStore, key string) {
	value, ok := kvs.Get(key)
	if ok {
		t.Errorf("Key '%s' has value '%s', expected that it doesn't exist", key, value)
	}
}

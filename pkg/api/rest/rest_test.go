package rest

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/istaheev/kvstore-assignment/pkg/kvstore"
	"github.com/stretchr/testify/assert"
)

func TestGetSet(t *testing.T) {
	kvs := kvstore.New()
	restHandler := New(&kvs)

	expectedValue := "first"

	r := httptest.NewRequest("POST", "http://example.com/key1", strings.NewReader(expectedValue))
	w := httptest.NewRecorder()
	restHandler.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)

	r = httptest.NewRequest("GET", "http://example.com/key1", nil)
	w = httptest.NewRecorder()
	restHandler.ServeHTTP(w, r)

	resp := w.Result()
	value, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, expectedValue, string(value))
}

func TestNonExistentKey(t *testing.T) {
	kvs := kvstore.New()
	restHandler := New(&kvs)

	r := httptest.NewRequest("GET", "http://example.com/key1", nil)
	w := httptest.NewRecorder()
	restHandler.ServeHTTP(w, r)

	assert.Equal(t, 404, w.Code)
}

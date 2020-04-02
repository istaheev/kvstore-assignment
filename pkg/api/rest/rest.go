package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/istaheev/kvstore-assignment/pkg/kvstore"
)

type keyGetHandler struct {
	kvs *kvstore.KeyValueStore
}

func (h keyGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	key := mux.Vars(r)["key"]
	value, ok := h.kvs.Get(key)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(value))
	}
}

type keySetHandler struct {
	kvs *kvstore.KeyValueStore
}

func (h keySetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	key := mux.Vars(r)["key"]
	valueArr, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Request body read error: %s", err.Error())))
		return
	}

	if !utf8.Valid(valueArr) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request body is not valid UTF-8"))
		return
	}

	value := string(valueArr)
	h.kvs.Set(key, value)
	w.Write([]byte("Key set"))
}

// New initializes REST-like API
func New(kvs *kvstore.KeyValueStore) http.Handler {
	var router = mux.NewRouter().StrictSlash(true)
	router.Handle("/{key}", keyGetHandler{kvs}).Methods(http.MethodGet)
	router.Handle("/{key}", keySetHandler{kvs}).Methods(http.MethodPost)
	return router
}

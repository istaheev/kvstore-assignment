package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/istaheev/kvstore-assignment/pkg/kvstore"
)

type keyGetHandler struct {
}

func (h keyGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var key = mux.Vars(r)["key"]
	w.Write([]byte(fmt.Sprintf("Requested key '%s'", key)))
}

type keySetHandler struct {
}

func (h keySetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var key = mux.Vars(r)["key"]
	w.Write([]byte(fmt.Sprintf("About to set key '%s'", key)))
}

// New initializes REST-like API
func New(kvs *kvstore.KeyValueStore) http.Handler {
	var router = mux.NewRouter().StrictSlash(true)
	router.Handle("/{key}", keyGetHandler{}).Methods(http.MethodGet)
	router.Handle("/{key}", keySetHandler{}).Methods(http.MethodPost)
	return router
}

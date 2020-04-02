package api

import (
	"net/http"

	"github.com/istaheev/kvstore-assignment/pkg/kvstore"

	"github.com/istaheev/kvstore-assignment/pkg/api/rest"
)

// NewServer creates new HTTP server with API handlers attached
func NewServer(listenAddr string, kvs *kvstore.KeyValueStore) *http.Server {
	return &http.Server{
		Addr:    listenAddr,
		Handler: logRequests(rest.New(kvs)),
	}
}

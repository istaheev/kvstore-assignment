package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/istaheev/kvstore-assignment/pkg/api"
	"github.com/istaheev/kvstore-assignment/pkg/kvstore"
)

//
// Command-line flags
//

var listenAddr = flag.String(
	"listen",
	":8080",
	"Address to listen for HTTP requests on",
)

// Remove expired keys from the store. Working on small batches we will give
// time to other requests to be processed. If there are no keys to remove then
// sleep a bit to not waste resources.
func expiredKeysRemovalLoop(kvs *kvstore.KeyValueStore) {
	for {
		removed := kvs.RemoveExpired(50)
		if removed > 0 {
			log.Printf("%d expired keys removed", removed)
		}
		if removed == 0 {
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	flag.Parse()

	// Initialize the store
	kvs := kvstore.New()
	kvs.SetExpiration(30 * time.Minute)
	go expiredKeysRemovalLoop(&kvs)

	// Run HTTP server to accept external requests
	log.Printf("Accepting HTTP requests on %s", *listenAddr)

	server := api.NewServer(*listenAddr, &kvs)

	// Explicit listener creation is required since default ListenAndServe() tends
	// to bind itself on tcp6 only
	l, err := net.Listen("tcp4", server.Addr)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Serve(l)
	if err != nil {
		log.Fatalln("Serve error:", err)
	}
}

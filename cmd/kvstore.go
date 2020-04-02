package main

import (
	"flag"
	"log"
	"net"

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

func main() {
	flag.Parse()

	kvs := kvstore.New()

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

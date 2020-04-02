package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

var server = flag.String(
	"server",
	"127.0.0.1:8080",
	"Address of the server to connect to",
)

var numReaders = 8
var numWriters = 2
var maxKey uint64 = 4294967296 // 2**32

var numWrites uint64 = 0
var numReads uint64 = 0

func chooseKey() string {
	key := rand.Uint64() % maxKey
	return fmt.Sprintf("%x", key)
}

func readAttempt() {
	key := chooseKey()
	resp, err := http.Get(fmt.Sprintf("http://%s/%s", *server, key))
	if err != nil {
		log.Println("http.Get error:", err)
		return
	}
	defer resp.Body.Close()

	resp = resp

	atomic.AddUint64(&numReads, 1)
}

func reader() {
	for {
		readAttempt()
	}
}

func writeAttempt() {
	key := chooseKey()
	resp, err := http.Post(
		fmt.Sprintf("http://%s/%s", *server, key),
		"plain/text",
		strings.NewReader("qwerty"),
	)
	if err != nil {
		log.Println("http.Get error:", err)
		return
	}
	defer resp.Body.Close()

	resp = resp

	atomic.AddUint64(&numWrites, 1)
}

func writer() {
	for {
		writeAttempt()
	}
}

func main() {
	flag.Parse()

	fmt.Println("Press CTRL-C to exit..")

	for i := 0; i < numReaders; i++ {
		go reader()
	}

	for i := 0; i < numWriters; i++ {
		go writer()
	}

	var prevNumReads uint64 = 0
	var prevNumWrites uint64 = 0
	var delay = 10 * time.Second

	for {
		time.Sleep(delay)

		reads := numReads - prevNumReads
		writes := numWrites - prevNumWrites

		fmt.Printf("Reads: %d (%.1f per sec); Writes: %d (%.1f per sec)\n",
			reads,
			float64(reads)/delay.Seconds(),
			writes,
			float64(writes)/delay.Seconds(),
		)
		fmt.Printf("TOTAL Reads: %d; Writes: %d\n\n", numReads, numWrites)

		prevNumReads = numReads
		prevNumWrites = numWrites
	}
}

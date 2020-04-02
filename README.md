# Downloading

    go get github.com/istaheev/kvstore-assignment

The application uses [dep](https://github.com/golang/dep) for dependency
management. Install the latest version via

    sudo apt-get install go-dep

or

    https://github.com/golang/dep

Then refresh dependencies:

    dep ensure

# Running

The entry point is located in `cmd/kvstore.go`:

    go run cmd/kvstore.go

The application accepts some command line arguments. Run

    go run cmd/kvstore.go -h

to get an overview of them.

There is a simple client with the main purpose to test the service under high
load:

    go run cmd/client/client.go

# Testing

    go test -v ./...



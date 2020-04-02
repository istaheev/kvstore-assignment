# Running

The entry point is located in `cmd/kvstore.go`:

    go run cmd/kvstore.go

The application accepts some command line arguments. Run

    go run cmd/kvstore.go -h

to get an overview of them.

# Building

The application uses [dep](https://github.com/golang/dep) for dependency management. Install the latest version via

    sudo apt-get install go-dep

or

    https://github.com/golang/dep

Then refresh dependencies:

    dep ensure


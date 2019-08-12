package main

import (
	"github.com/eqlabs/sprawl/api"
	"github.com/eqlabs/sprawl/db"
)

func main() {
	// Start up the database
	storage := &db.Storage{}
	storage.SetDbPath("/var/lib/sprawl/data")
	storage.Run()
	defer storage.Close()

	// p2pImpl := &p2p.P2p{}
	api.Run(storage, 1337)
}

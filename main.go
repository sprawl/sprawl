package main

import (
	"github.com/eqlabs/sprawl/api"
	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/db"
	"github.com/eqlabs/sprawl/p2p"
)

func main() {
	// Load config
	config := &config.Config{}
	config.ReadConfig("config/default")

	// Start up the database
	storage := &db.Storage{}
	storage.SetDbPath(config.GetString("database.path"))
	storage.Run()
	defer storage.Close()

	// Run the P2P process
	p2pInstance := p2p.P2p{}
	p2pInstance.Run()

	// Run the gRPC API
	api.Run(storage, config.GetUint("api.port"))
}

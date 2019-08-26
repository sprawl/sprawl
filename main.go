package main

import (
	"fmt"

	"github.com/eqlabs/sprawl/api"
	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/db"
	"github.com/eqlabs/sprawl/p2p"
)

func main() {
	// Load config
	config := &config.Config{}
	config.ReadConfig("config/default")

	fmt.Println(config.GetString("database.path"))

	// Start up the database
	storage := &db.Storage{}
	storage.SetDbPath(config.GetString("database.path"))
	storage.Run()
	defer storage.Close()

	// Run the P2P process
	p2pInstance := p2p.NewP2p()
	p2pInstance.Run()

	// Run the gRPC API
	api.Run(storage, config.GetUint("api.port"))
}

package main

import (
	"fmt"

	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/db"
	"github.com/eqlabs/sprawl/p2p"
	"github.com/eqlabs/sprawl/service"
)

func main() {
	// Load config
	config := &config.Config{}
	config.ReadConfig("config/default")

	fmt.Printf("Saving data to %s", config.GetString("database.path"))

	// Start up the database
	storage := &db.Storage{}
	storage.SetDbPath(config.GetString("database.path"))
	storage.Run()
	defer storage.Close()

	// Run the P2P process
	p2pInstance := p2p.NewP2p()
	p2pInstance.Run()

	// Construct the server struct
	server := service.Server{}
	p2pInstance.RegisterOrderService(&server.Orders)
	p2pInstance.RegisterChannelService(&server.Channels)

	// Run the gRPC API
	server.Run(storage, p2pInstance, config.GetUint("api.port"))
}

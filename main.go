package main

import (
	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/db"
	"github.com/eqlabs/sprawl/interfaces"
	"github.com/eqlabs/sprawl/p2p"
	"github.com/eqlabs/sprawl/service"
	"github.com/prometheus/common/log"
)

func main() {
	// Load config
	var config interfaces.Config = &config.Config{}
	config.ReadConfig("config/default")

	log.Infof("Saving data to %s", config.GetString("database.path"))

	// Start up the database
	var storage interfaces.Storage = &db.Storage{}
	storage.SetDbPath(config.GetString("database.path"))
	storage.Run()
	defer storage.Close()

	// Run the P2P process
	p2pInstance := p2p.NewP2p()
	p2pInstance.Run()

	// Construct the server struct
	server := service.NewServer(storage, p2pInstance)

	// Connect the order and channel services with p2p
	p2pInstance.RegisterOrderService(server.Orders)
	p2pInstance.RegisterChannelService(server.Channels)

	// Run the gRPC API
	server.Run(config.GetUint("api.port"))
}

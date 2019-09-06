package main

import (
	"time"

	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/db"
	"github.com/eqlabs/sprawl/interfaces"
	"github.com/eqlabs/sprawl/p2p"
	"github.com/eqlabs/sprawl/pb"
	"github.com/eqlabs/sprawl/service"
	"github.com/gogo/protobuf/proto"
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
	if config.GetString("api.debug_pinger") == "true" {
		debugPinger(p2pInstance)
	}

	// Construct the server struct
	server := service.NewServer(storage, p2pInstance)

	// Connect the order and channel services with p2p
	p2pInstance.RegisterOrderService(server.Orders)
	p2pInstance.RegisterChannelService(server.Channels)

	// Run the gRPC API
	server.Run(config.GetUint("api.port"))
}

func debugPinger(p2pInstance *p2p.P2p) {
	var testChannel *pb.Channel = &pb.Channel{Id: []byte("testChannel")}
	var testOrder *pb.Order = &pb.Order{Asset: string("ETH"), CounterAsset: string("BTC"), Amount: 52152, Price: 0.2, Id: []byte("jgkahgkjal")}
	testOrderInBytes, err := proto.Marshal(testOrder)
	if err != nil {
		panic(err)
	}
	testWireMessage := &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}
	for {
		log.Infof("Debug pinger is sending testWireMessage: %s\n", testWireMessage)
		p2pInstance.Send(testWireMessage)
		time.Sleep(time.Minute)
	}
}

package service

import (
	fmt "fmt"
	"log"
	"net"

	"github.com/eqlabs/sprawl/interfaces"
	"github.com/eqlabs/sprawl/pb"
	"google.golang.org/grpc"
)

// Server contains services for both Orders and Channels
type Server struct {
	Orders   OrderService
	Channels ChannelService
}

// Run runs the gRPC server
func (server Server) Run(storage interfaces.Storage, p2p interfaces.P2p, port uint) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	// Create an OrderService that defines the order handling operations
	server.Orders = OrderService{}
	server.Orders.RegisterStorage(storage)
	server.Orders.RegisterP2p(p2p)

	// Create an ChannelService that defines channel operations
	server.Channels = ChannelService{}
	server.Channels.RegisterStorage(storage)
	server.Channels.RegisterP2p(p2p)

	// Register the Services with the RPC server
	pb.RegisterOrderHandlerServer(s, &server.Orders)
	pb.RegisterChannelHandlerServer(s, &server.Channels)

	// Run the server
	s.Serve(lis)
}

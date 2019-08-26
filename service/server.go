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
	// Listen to TCP connections
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Set options, here we can configure things like TLS support
	opts := []grpc.ServerOption{}
	// Create new gRPC server with (blank) options
	s := grpc.NewServer(opts...)

	// Create an OrderService that defines the order handling operations
	server.Orders = OrderService{}
	server.Orders.RegisterStorage(storage)
	server.Orders.RegisterP2p(p2p)

	// Create an OrderService that stores the endpoints
	server.Channels = ChannelService{}
	server.Channels.RegisterStorage(storage)
	server.Channels.RegisterP2p(p2p)

	// Register the Services with the RPC server
	pb.RegisterOrderHandlerServer(s, &server.Orders)
	pb.RegisterChannelHandlerServer(s, &server.Channels)

	// Run the server
	s.Serve(lis)
}

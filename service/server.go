package service

import (
	fmt "fmt"
	"net"

	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/interfaces"
	"github.com/sprawl/sprawl/pb"
	"google.golang.org/grpc"
)

// Server contains services for both Orders and Channels
type Server struct {
	Orders   *OrderService
	Channels *ChannelService
	Logger   interfaces.Logger
}

// NewServer returns a server that has connections to p2p and storage
func NewServer(log interfaces.Logger, storage interfaces.Storage, p2p interfaces.P2p) *Server {
	server := &Server{Logger: log}

	// Create an OrderService that defines the order handling operations
	server.Orders = &OrderService{Logger: log}
	server.Orders.RegisterStorage(storage)
	server.Orders.RegisterP2p(p2p)

	// Create an ChannelService that defines channel operations
	server.Channels = &ChannelService{}
	server.Channels.RegisterStorage(storage)
	server.Channels.RegisterP2p(p2p)

	return server
}

// Run runs the gRPC server
func (server *Server) Run(port uint) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if !errors.IsEmpty(err) {
		if server.Logger != nil {
			server.Logger.Fatal(errors.E(errors.Op("Listen"), err))
		}
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	// Register the Services with the RPC server
	pb.RegisterOrderHandlerServer(s, server.Orders)
	pb.RegisterChannelHandlerServer(s, server.Channels)

	// Run the server
	s.Serve(lis)
}

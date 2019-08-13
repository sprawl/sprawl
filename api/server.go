package api

import (
	fmt "fmt"
	"log"
	"net"

	"github.com/eqlabs/sprawl/db"
	"google.golang.org/grpc"
)

// Run runs the gRPC server
func Run(storage *db.Storage, port uint) {
	// Listen to TCP connections
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Set options, here we can configure things like TLS support
	opts := []grpc.ServerOption{}
	// Create new gRPC server with (blank) options
	s := grpc.NewServer(opts...)

	// Create an OrderService that stores the endpoints
	service := &OrderService{}
	// Register the storage service with it
	service.RegisterStorage(storage)

	// Register the OrderService with the server
	RegisterOrderHandlerServer(s, service)
	// Run the server
	s.Serve(lis)
}

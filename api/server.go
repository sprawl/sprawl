package api

import (
	fmt "fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

// Init initializes the gRPC server
func Init() {
	// Listen to TCP connections
	port := 1337
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Set options, here we can configure things like TLS support
	opts := []grpc.ServerOption{}
	// Create new gRPC server with (blank) options
	s := grpc.NewServer(opts...)
	// Register the service with the server
	RegisterOrderHandlerServer(s, &OrderService{})
	// Run the server
	s.Serve(lis)
}

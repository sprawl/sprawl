package api

import (
	fmt "fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

// Run runs the gRPC server
func Run(port uint16) {
	// Listen to TCP connections
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

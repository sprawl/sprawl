package api

import (
	"context"
	fmt "fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

// OrderService is a implementation of OrderService Grpc Service.
type OrderService struct{}

// Create function implementation of gRPC Service.
func (s *OrderService) Create(ctx context.Context, in *Order) (*CreateResponse, error) {
	return &CreateResponse{
		CreatedOrder: in,
		Error:        nil,
	}, nil
}

// Init initializes the grpc server
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

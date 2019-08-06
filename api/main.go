package api

import (
	"context"

	"google.golang.org/grpc"
)

// OrderServiceImpl is a implementation of OrderService Grpc Service.
type OrderServiceImpl struct{}

// NewOrderServiceImpl returns the pointer to the implementation.
func NewOrderServiceImpl() *OrderServiceImpl {
	return &grpc.ServiceDesc{}
}

// Create function implementation of gRPC Service.
func (serviceImpl *OrderServiceImpl) Create(ctx context.Context, order *Order) (*CreateResponse, error) {
	return &CreateResponse{
		CreatedOrder: order,
		Error:        nil,
	}, nil
}

// Init initializes the grpc server
func Init() {
	order := Order{
		id: 1234,
	}
	// Set options, here we can configure things like TLS support
	opts := []grpc.ServerOption{}
	// Create new gRPC server with (blank) options
	s := grpc.NewServer(opts...)
	// Create BlogService type
	orderService := &OrderServiceImpl{}
	// Register the service with the server
	s.RegisterService(orderService)
}

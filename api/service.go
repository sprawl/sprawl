package api

import "context"

// OrderService implements the OrderService Server service.proto
type OrderService struct{}

// Create function implementation of gRPC Service.
func (s *OrderService) Create(ctx context.Context, in *Order) (*CreateResponse, error) {
	return &CreateResponse{
		CreatedOrder: in,
		Error:        nil,
	}, nil
}

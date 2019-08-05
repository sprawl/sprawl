package api

import (
	"google.golang.org/grpc"
	"context"
)

func Init {
	order := pb.Order {
		id:    1234,
	}
}

//OrderServiceImpl is a implementation of OrderService Grpc Service.
type OrderServiceImpl struct {}

//NewOrderServiceImpl returns the pointer to the implementation.
func NewOrderServiceImpl() *OrderServiceImpl {
	return &OrderServiceImpl{}
}

//Add function implementation of gRPC Service.
func (serviceImpl *OrderServiceImpl) Create(ctx context.Context, order *Order) (*CreateResponse, error) {
	return &CreateResponse {
		CreatedOrder: order,
		Error:           nil,
		}, nil
}

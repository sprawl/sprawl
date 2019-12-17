package interfaces

import (
	"context"

	"github.com/sprawl/sprawl/pb"
)

// OrderService is an interface to the Order endpoints in sprawl.proto
type OrderService interface {
	RegisterStorage(db Storage)
	RegisterP2p(p2p P2p)
	Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error)
	Receive(data []byte) error
	Delete(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.GenericResponse, error)
	Lock(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.GenericResponse, error)
	Unlock(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.GenericResponse, error)
	GetOrder(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.Order, error)
	GetAllOrders(ctx context.Context, in *pb.Empty) (*pb.OrderList, error)
}

package interfaces

import (
	"context"

	"github.com/eqlabs/sprawl/pb"
)

type OrderService interface {
	RegisterStorage(db Storage)
	RegisterP2p(p2p P2p)
	Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error)
	Delete(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.GenericResponse, error)
	Lock(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.GenericResponse, error)
	Unlock(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.GenericResponse, error)
}

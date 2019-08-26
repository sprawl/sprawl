package interfaces

import (
	"context"

	"github.com/eqlabs/sprawl/pb"
)

type ChannelService interface {
	RegisterStorage(db Storage)
	RegisterP2p(p2p P2p)
	Join(ctx context.Context, in *pb.Channel) (*pb.JoinResponse, error)
	Leave(ctx context.Context, in *pb.Channel) (*pb.GenericResponse, error)
}

package interfaces

import (
	"context"

	"github.com/eqlabs/sprawl/pb"
)

// ChannelService is an interface to the Channel endpoints in sprawl.proto
type ChannelService interface {
	RegisterStorage(db Storage)
	RegisterP2p(p2p P2p)
	Join(ctx context.Context, in *pb.ChannelOptions) (*pb.JoinResponse, error)
	Leave(ctx context.Context, in *pb.Channel) (*pb.GenericResponse, error)
	GetChannel(ctx context.Context, in *pb.ChannelSpecificRequest) (*pb.Channel, error)
	GetAllChannels(ctx context.Context, in *pb.Empty) (*pb.ChannelListResponse, error)
}

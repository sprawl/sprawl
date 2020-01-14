package interfaces

import (
	"context"

	"github.com/sprawl/sprawl/pb"
)

type NodeService interface {
	RegisterP2p(p2p P2p)
	GetAllPeers(ctx context.Context, in *pb.Empty) (*pb.PeerListResponse, error)
	BlacklistPeer(ctx context.Context, in *pb.Peer) (*pb.GenericResponse, error)
}

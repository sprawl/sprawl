package service

import (
	"context"

	"github.com/sprawl/sprawl/interfaces"
	"github.com/sprawl/sprawl/pb"
)

// ChannelService implements the ChannelHandlerServer service.proto
type NodeService struct {
	P2p interfaces.P2p
}

func (s *NodeService) RegisterP2p(p2p interfaces.P2p) {
	s.P2p = p2p
}

func (s *NodeService) GetAllPeers(ctx context.Context, in *pb.Empty) (*pb.PeerListResponse, error) {
	data := s.P2p.GetAllPeers()
	peerList := &pb.PeerListResponse{PeerIds: data}
	return peerList, nil
}

func (s *NodeService) BlacklistPeer(ctx context.Context, in *pb.PeerId) (*pb.GenericResponse, error) {
	s.P2p.BlacklistPeer(in)
	return &pb.GenericResponse{
		Error: nil,
	}, nil
}

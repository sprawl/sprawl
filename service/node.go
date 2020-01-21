package service

import (
	"context"

	"github.com/sprawl/sprawl/interfaces"
	"github.com/sprawl/sprawl/pb"
)

// NodeService is a gRPC service for p2p operations.
type NodeService struct {
	P2p interfaces.P2p
}

// RegisterP2p registers a p2p interface with NodeService
func (s *NodeService) RegisterP2p(p2p interfaces.P2p) {
	s.P2p = p2p
}

// GetAllPeers fetches all connected peers from NodeService.P2p
func (s *NodeService) GetAllPeers(ctx context.Context, in *pb.Empty) (*pb.PeerListResponse, error) {
	data := s.P2p.GetAllPeers()
	peerList := &pb.PeerListResponse{PeerIds: data}
	return peerList, nil
}

// BlacklistPeer blacklists a peer from connecting to this node
func (s *NodeService) BlacklistPeer(ctx context.Context, in *pb.Peer) (*pb.Empty, error) {
	s.P2p.BlacklistPeer(in)
	return &pb.Empty{}, nil
}

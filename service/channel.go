package service

import (
	"context"

	"github.com/eqlabs/sprawl/interfaces"
	"github.com/eqlabs/sprawl/pb"
)

// ChannelService implements the ChannelService Server service.proto
type ChannelService struct {
	storage interfaces.Storage
	p2p     interfaces.P2p
}

// RegisterStorage registers a storage service to store the Channels in
func (s *ChannelService) RegisterStorage(storage interfaces.Storage) {
	s.storage = storage
}

// RegisterP2p registers a p2p service
func (s *ChannelService) RegisterP2p(p2p interfaces.P2p) {
	s.p2p = p2p
}

// Join joins a channel, subscribing to new topic in libp2p
func (s *ChannelService) Join(ctx context.Context, in *pb.Channel) (*pb.JoinResponse, error) {
	channelID := in.GetId()

	s.p2p.Subscribe(string(channelID))

	return &pb.JoinResponse{
		JoinedChannel: &pb.Channel{Id: channelID},
	}, nil
}

// Leave leaves a channel, removing a subscription from libp2p
func (s *ChannelService) Leave(ctx context.Context, in *pb.Channel) (*pb.GenericResponse, error) {

	// TODO: Add Channel leaving logic

	return &pb.GenericResponse{
		Error: nil,
	}, nil
}

package service

import (
	"context"
	"sort"
	"strings"

	"github.com/eqlabs/sprawl/interfaces"
	"github.com/eqlabs/sprawl/pb"
)

// ChannelService implements the ChannelHandlerServer service.proto
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
func (s *ChannelService) Join(ctx context.Context, in *pb.JoinRequest) (*pb.JoinResponse, error) {
	// Get all channel options, sort
	assetPair := []string{string(in.GetAsset()), string(in.GetCounterAsset())}
	sort.Strings(assetPair)

	// Join the channel options together
	channelOptBlob := []byte(strings.Join(assetPair[:], ","))

	s.p2p.Subscribe(string(channelOptBlob))

	joinedChannel := &pb.Channel{Id: channelOptBlob}

	return &pb.JoinResponse{
		JoinedChannel: joinedChannel,
	}, nil
}

// Leave leaves a channel, removing a subscription from libp2p
func (s *ChannelService) Leave(ctx context.Context, in *pb.Channel) (*pb.GenericResponse, error) {

	// TODO: Add Channel leaving logic

	return &pb.GenericResponse{
		Error: nil,
	}, nil
}

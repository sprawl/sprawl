package service

import (
	"context"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/interfaces"
	"github.com/sprawl/sprawl/pb"
)

// ChannelService implements the ChannelHandlerServer service.proto
type ChannelService struct {
	Storage interfaces.Storage
	P2p     interfaces.P2p
}

func getChannelStorageKey(channelOptBlob []byte) []byte {
	return []byte(strings.Join([]string{string(interfaces.ChannelPrefix), string(channelOptBlob)}, ""))
}

// RegisterStorage registers a storage service to store the Channels in
func (s *ChannelService) RegisterStorage(storage interfaces.Storage) {
	s.Storage = storage
}

// RegisterP2p registers a p2p service
func (s *ChannelService) RegisterP2p(p2p interfaces.P2p) {
	s.P2p = p2p
}

// Join joins a channel, subscribing to new topic in libp2p
func (s *ChannelService) Join(ctx context.Context, in *pb.JoinRequest) (*pb.JoinResponse, error) {
	// Get all channel options, sort
	assetPair := []string{string(in.GetAsset()), string(in.GetCounterAsset())}
	sort.Strings(assetPair)

	// Join the channel options together
	channelOptBlob := []byte(strings.Join(assetPair[:], ","))

	// Create a Channel protobuf message to return to the user
	joinedChannel := &pb.Channel{Id: channelOptBlob, Options: &pb.ChannelOptions{AssetPair: strings.Join(assetPair, "")}}
	marshaledChannel, err := proto.Marshal(joinedChannel)
	if !errors.IsEmpty(err) {
		return nil, errors.E(errors.Op("Join"), err)
	}

	// Subscribe to a topic matching the options
	s.P2p.Subscribe(joinedChannel)

	// Store the joined channel in LevelDB
	s.Storage.Put(getChannelStorageKey(channelOptBlob), marshaledChannel)

	return &pb.JoinResponse{
		JoinedChannel: joinedChannel,
	}, nil
}

// Leave leaves a channel, removing a subscription from libp2p
func (s *ChannelService) Leave(ctx context.Context, in *pb.ChannelSpecificRequest) (*pb.GenericResponse, error) {
	channelOptBlob := in.GetId()

	// Remove the channel from LevelDB
	s.Storage.Delete(getChannelStorageKey(channelOptBlob))

	return &pb.GenericResponse{
		Error: nil,
	}, nil
}

// GetChannel fetches a single channel from the database
func (s *ChannelService) GetChannel(ctx context.Context, in *pb.ChannelSpecificRequest) (*pb.Channel, error) {
	data, err := s.Storage.Get(getChannelStorageKey(in.GetId()))
	if !errors.IsEmpty(err) {
		return nil, errors.E(errors.Op("Get channel"), err)
	}
	channel := &pb.Channel{}
	proto.Unmarshal(data, channel)
	return channel, nil
}

// GetAllChannels fetches all channels from the database
func (s *ChannelService) GetAllChannels(ctx context.Context, in *pb.Empty) (*pb.ChannelList, error) {
	data, err := s.Storage.GetAllWithPrefix(string(interfaces.ChannelPrefix))
	if !errors.IsEmpty(err) {
		return nil, errors.E(errors.Op("Get all channels "), err)
	}

	channels := make([]*pb.Channel, 0)
	i := 0
	for _, value := range data {
		channel := &pb.Channel{}
		proto.Unmarshal([]byte(value), channel)
		channels = append(channels, channel)
		i++
	}

	ChannelList := &pb.ChannelList{Channels: channels}
	return ChannelList, nil
}

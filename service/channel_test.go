package service

import (
	"testing"

	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/interfaces"
	"github.com/sprawl/sprawl/pb"
	"github.com/stretchr/testify/assert"
)

func leaveTestChannel() {
	p2pInstance.Unsubscribe(channel)
}

func TestChannelStorageKeyPrefixer(t *testing.T) {
	prefixedBytes := getChannelStorageKey([]byte(asset1))
	assert.Equal(t, string(prefixedBytes), string(interfaces.ChannelPrefix)+asset1)
}

func TestChannelJoining(t *testing.T) {
	leaveTestChannel()

	createNewServerInstance()
	defer p2pInstance.Close()
	defer storage.Close()
	defer conn.Close()

	// Create a ChannelService that stores the endpoints
	var channelService interfaces.ChannelService = &ChannelService{}
	// Register the services
	channelService.RegisterStorage(storage)
	channelService.RegisterP2p(p2pInstance)

	var lastChannel *pb.Channel

	// Register channel endpoints with the gRPC server
	pb.RegisterChannelHandlerServer(s, channelService)

	go func() {
		if err := s.Serve(lis); !errors.IsEmpty(err) {
			log.Fatalf("Server exited with error: %v", err)
		}
		defer s.Stop()
	}()

	var channelClient pb.ChannelHandlerClient = pb.NewChannelHandlerClient(conn)

	resp, err := channelClient.Join(ctx, &pb.JoinRequest{Asset: asset1, CounterAsset: asset2})
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	_, err = channelClient.Join(ctx, &pb.JoinRequest{Asset: asset2, CounterAsset: asset1})
	assert.Error(t, err)

	lastChannel = resp.GetJoinedChannel()
	storedChannel, err := channelClient.GetChannel(ctx, &pb.ChannelSpecificRequest{Id: lastChannel.GetId()})
	assert.NoError(t, err)
	assert.Equal(t, lastChannel, storedChannel)

	resp3, err := channelClient.GetAllChannels(ctx, &pb.Empty{})
	channelList := resp3.GetChannels()
	assert.Equal(t, 1, len(channelList))

	_, err = channelClient.Leave(ctx, &pb.ChannelSpecificRequest{Id: lastChannel.GetId()})
	assert.NoError(t, err)
}

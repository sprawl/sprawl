package service

import (
	"testing"

	"github.com/eqlabs/sprawl/interfaces"
	"github.com/eqlabs/sprawl/pb"
	"github.com/stretchr/testify/assert"
)

func leaveEveryChannel() {
	storage.DeleteAllWithPrefix(string(interfaces.ChannelPrefix))
}

func TestChannelStorageKeyPrefixer(t *testing.T) {
	prefixedBytes := getChannelStorageKey([]byte(asset1))
	assert.Equal(t, string(prefixedBytes), string(interfaces.ChannelPrefix)+asset1)
}

func TestChannelJoining(t *testing.T) {
	createNewServerInstance()
	defer p2pInstance.Close()
	defer storage.Close()
	defer conn.Close()
	leaveEveryChannel()

	// Create a ChannelService that stores the endpoints
	var channelService interfaces.ChannelService = &ChannelService{}
	// Register the services
	channelService.RegisterStorage(storage)
	channelService.RegisterP2p(p2pInstance)

	var lastChannel *pb.Channel

	// Register channel endpoints with the gRPC server
	pb.RegisterChannelHandlerServer(s, channelService)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
		defer s.Stop()
	}()

	var channelClient pb.ChannelHandlerClient = pb.NewChannelHandlerClient(conn)

	resp, err := channelClient.Join(ctx, &pb.JoinRequest{Asset: asset1, CounterAsset: asset2})
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp2, err := channelClient.Join(ctx, &pb.JoinRequest{Asset: asset2, CounterAsset: asset1})
	assert.NoError(t, err)
	assert.NotNil(t, resp2)

	assert.Equal(t, resp.GetJoinedChannel().GetId(), resp2.GetJoinedChannel().GetId())

	lastChannel = resp.GetJoinedChannel()
	t.Logf("Last channel: %s", lastChannel)

	storedChannel, err := channelClient.GetChannel(ctx, &pb.ChannelSpecificRequest{Id: lastChannel.GetId()})
	assert.NoError(t, err)
	assert.Equal(t, lastChannel, storedChannel)

	resp3, err := channelClient.GetAllChannels(ctx, &pb.Empty{})
	channelList := resp3.GetChannels()
	assert.Equal(t, len(channelList), 1)

	_, err = channelClient.Leave(ctx, &pb.ChannelSpecificRequest{Id: lastChannel.GetId()})
	assert.NoError(t, err)
}

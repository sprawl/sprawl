package service

import (
	"log"
	"testing"

	"github.com/eqlabs/sprawl/interfaces"
	"github.com/eqlabs/sprawl/pb"
	"github.com/stretchr/testify/assert"
)

func TestChannelStorageKeyPrefixer(t *testing.T) {
	prefixedBytes := getChannelStorageKey([]byte(asset1))
	assert.Equal(t, string(prefixedBytes), string(interfaces.ChannelPrefix)+string(asset1))
}

func TestChannelJoining(t *testing.T) {
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
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
		defer s.Stop()
	}()

	var channelClient pb.ChannelHandlerClient = pb.NewChannelHandlerClient(conn)

	resp, err := channelClient.Join(ctx, &pb.ChannelOptions{Asset: []byte(asset1), CounterAsset: []byte(asset2)})
	assert.Equal(t, err, nil)
	assert.NotEqual(t, resp, nil)

	lastChannel = resp.GetJoinedChannel()
	t.Log(lastChannel)
	storedChannel, err := channelClient.GetChannel(ctx, &pb.ChannelSpecificRequest{Id: lastChannel.GetId()})
	assert.Equal(t, err, nil)
	assert.Equal(t, lastChannel, storedChannel)
}

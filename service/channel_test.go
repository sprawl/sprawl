package service

import (
	"context"
	"log"
	"testing"

	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/db"
	"github.com/eqlabs/sprawl/interfaces"
	"github.com/eqlabs/sprawl/p2p"
	"github.com/eqlabs/sprawl/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	bufconn "google.golang.org/grpc/test/bufconn"
)

func TestChannelStorageKeyPrefixer(t *testing.T) {
	prefixedBytes := getChannelStorageKey([]byte(asset1))
	assert.Equal(t, string(prefixedBytes), string(interfaces.ChannelPrefix)+string(asset1))
}

func TestChannelJoining(t *testing.T) {
	var storage *db.Storage = &db.Storage{}
	var p2pInstance *p2p.P2p = p2p.NewP2p()
	p2pInstance.Run()
	defer p2pInstance.Close()

	// Load config
	config := &config.Config{}
	config.ReadConfig(testConfigPath)

	// Initialize storage
	storage.SetDbPath(config.GetString(dbPathVar))
	storage.Run()
	defer storage.Close()

	ctx = context.Background()
	lis = bufconn.Listen(bufSize)

	conn, err = grpc.DialContext(ctx, dialContext, grpc.WithDialer(BufDialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Create gRPC server
	s := grpc.NewServer()

	// Create a ChannelService that stores the endpoints
	var channelService interfaces.ChannelService = &ChannelService{}
	// Register the services
	channelService.RegisterStorage(storage)
	channelService.RegisterP2p(p2pInstance)

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
}

package service

import (
	"context"
	"log"
	"testing"

	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/db"
	"github.com/eqlabs/sprawl/interfaces"
	"github.com/eqlabs/sprawl/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	bufconn "google.golang.org/grpc/test/bufconn"
)

func TestChannelJoining(t *testing.T) {
	bufSize := 1024 * 1024

	// Load config
	var config interfaces.Config = &config.Config{}
	config.ReadConfig("../config/test")

	// Initialize storage
	var storage interfaces.Storage = &db.Storage{}
	storage.SetDbPath(config.GetString("database.path"))
	storage.Run()

	// "Listen" to buffer
	lis := bufconn.Listen(bufSize)

	// Create gRPC server
	s := grpc.NewServer()

	// Create a ChannelService that stores the endpoints
	var channelService interfaces.ChannelService = &ChannelService{}

	// Register the storage service with it
	channelService.RegisterStorage(storage)

	pb.RegisterChannelHandlerServer(s, channelService)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(BufDialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	var channelClient pb.ChannelHandlerClient = pb.NewChannelHandlerClient(conn)

	resp, err := channelClient.Join(ctx, &pb.Channel{Id: []byte("arbitrarychannel")})
	assert.Equal(t, nil, err)
	t.Log("Joined Channel: ", resp)
	assert.NotEqual(t, false, resp)
}

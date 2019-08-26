package service

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/db"
	"github.com/eqlabs/sprawl/interfaces"
	"github.com/eqlabs/sprawl/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	bufconn "google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener
var orderClient pb.OrderHandlerClient
var channelClient pb.ChannelHandlerClient
var ctx context.Context

var storage *db.Storage

// Construct a correct test Order
var testOrder = pb.CreateRequest{Asset: []byte("ETH"), CounterAsset: []byte("BTC"), Amount: 52617562718, Price: 0.1}
var lastOrder *pb.Order

func init() {
	// Load config
	var config interfaces.Config = &config.Config{}
	config.ReadConfig("../config/test")

	// Initialize storage
	var storage interfaces.Storage = &db.Storage{}
	storage.SetDbPath(config.GetString("database.path"))
	storage.Run()

	// "Listen" to buffer
	lis = bufconn.Listen(bufSize)

	// Create gRPC server
	s := grpc.NewServer()

	// Create an OrderService that stores the endpoints
	var orderService interfaces.OrderService = &OrderService{}
	//var channelService interfaces.ChannelService = ChannelService{}

	// Register the storage service with it
	orderService.RegisterStorage(storage)

	pb.RegisterOrderHandlerServer(s, orderService)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	ctx = context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	orderClient = pb.NewOrderHandlerClient(conn)
	channelClient = pb.NewChannelHandlerClient(conn)
}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func TestOrderCreation(t *testing.T) {
	resp, err := orderClient.Create(ctx, &testOrder)
	assert.Equal(t, nil, err)
	t.Log("Created Order: ", resp)
	assert.NotEqual(t, false, resp)

	lastOrder = resp.GetCreatedOrder()

	resp2, err := orderClient.Delete(ctx, &pb.OrderSpecificRequest{Id: lastOrder.GetId()})
	assert.Equal(t, nil, err)
	assert.NotEqual(t, false, resp2)
}

/* func TestChannelJoining(t *testing.T) {
	resp, err := channelClient.Join(ctx, &pb.Channel{})
	assert.Equal(t, nil, err)
	t.Log("Joined Channel: ", resp)
	assert.NotEqual(t, false, resp)
} */

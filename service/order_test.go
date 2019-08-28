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
	"github.com/eqlabs/sprawl/p2p"
	"github.com/eqlabs/sprawl/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	bufconn "google.golang.org/grpc/test/bufconn"
)

var bufSize = 1024 * 1024
var lis *bufconn.Listener
var conn *grpc.ClientConn
var err error
var ctx context.Context

func BufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func TestOrderCreation(t *testing.T) {
	var storage *db.Storage = &db.Storage{}
	var p2pInstance *p2p.P2p = p2p.NewP2p()
	p2pInstance.Run()
	defer p2pInstance.Close()

	// Load config
	config := &config.Config{}
	config.ReadConfig("../config/test")

	// Initialize storage
	storage.SetDbPath(config.GetString("database.path"))
	storage.Run()
	defer storage.Close()

	ctx = context.Background()
	lis = bufconn.Listen(bufSize)

	conn, err = grpc.DialContext(ctx, "OrderEndpoint", grpc.WithDialer(BufDialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Create gRPC server
	s := grpc.NewServer()

	testOrder := pb.CreateRequest{Asset: []byte("ETH"), CounterAsset: []byte("BTC"), Amount: 52617562718, Price: 0.1}

	var lastOrder *pb.Order

	// Create an OrderService
	var orderService interfaces.OrderService = &OrderService{}
	// Register services
	orderService.RegisterStorage(storage)
	orderService.RegisterP2p(p2pInstance)
	// Register order endpoints with the gRPC server
	pb.RegisterOrderHandlerServer(s, orderService)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
		defer s.Stop()
	}()

	var orderClient pb.OrderHandlerClient = pb.NewOrderHandlerClient(conn)

	resp, err := orderClient.Create(ctx, &testOrder)
	assert.Equal(t, nil, err)
	t.Log("Created Order: ", resp)
	assert.NotEqual(t, false, resp)

	lastOrder = resp.GetCreatedOrder()

	resp2, err := orderClient.Delete(ctx, &pb.OrderSpecificRequest{Id: lastOrder.GetId()})
	assert.Equal(t, nil, err)
	assert.NotEqual(t, false, resp2)
}

func TestChannelJoining(t *testing.T) {
	var storage *db.Storage = &db.Storage{}
	var p2pInstance *p2p.P2p = p2p.NewP2p()
	p2pInstance.Run()
	defer p2pInstance.Close()

	// Load config
	config := &config.Config{}
	config.ReadConfig("../config/test")

	// Initialize storage
	storage.SetDbPath(config.GetString("database.path"))
	storage.Run()
	defer storage.Close()

	ctx = context.Background()
	lis = bufconn.Listen(bufSize)

	conn, err = grpc.DialContext(ctx, "ChannelEndpoint", grpc.WithDialer(BufDialer), grpc.WithInsecure())
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

	resp, err := channelClient.Join(ctx, &pb.JoinRequest{Asset: []byte("ETH"), CounterAsset: []byte("BTC")})

	assert.Equal(t, err, nil)
	assert.NotEqual(t, resp, nil)
}

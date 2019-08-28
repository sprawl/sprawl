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

const testConfigPath = "../config/test"
const dbPathVar = "database.path"
const dialContext = "TestEndpoint"
const asset1 = "ETH"
const asset2 = "BTC"
const testAmount = 52617562718
const testPrice = 0.1

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

	testOrder := pb.CreateRequest{Asset: []byte(asset1), CounterAsset: []byte(asset2), Amount: testAmount, Price: testPrice}

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

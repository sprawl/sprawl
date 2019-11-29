package service

import (
	"context"
	"crypto/rand"
	"net"
	"testing"
	"time"

	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/config"
	"github.com/sprawl/sprawl/db"
	"github.com/sprawl/sprawl/identity"
	"github.com/sprawl/sprawl/interfaces"
	"github.com/sprawl/sprawl/p2p"
	"github.com/sprawl/sprawl/pb"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	bufconn "google.golang.org/grpc/test/bufconn"
)

const testConfigPath string = "../config/test"
const dbPathVar string = "database.path"
const dialContext string = "TestEndpoint"
const asset1 string = "ETH"
const asset2 string = "BTC"
const testAmount = 52617562718
const testPrice = 0.1

var bufSize = 1024 * 1024
var lis *bufconn.Listener
var conn *grpc.ClientConn
var err error
var ctx context.Context
var storage *db.Storage = &db.Storage{}
var p2pInstance *p2p.P2p
var testConfig *config.Config
var s *grpc.Server
var orderClient pb.OrderHandlerClient
var orderService interfaces.OrderService = &OrderService{}
var channelService interfaces.ChannelService = &ChannelService{}
var channel *pb.Channel
var logger *zap.Logger
var log *zap.SugaredLogger

func init() {
	logger, _ = zap.NewProduction()
	log = logger.Sugar()
	testConfig = &config.Config{Logger: log}
	privateKey, publicKey, _ := identity.GenerateKeyPair(rand.Reader)
	p2pInstance = p2p.NewP2p(log, testConfig, privateKey, publicKey)
	testConfig.ReadConfig(testConfigPath)
	storage.SetDbPath(testConfig.GetString(dbPathVar))
}

func createNewServerInstance() {
	p2pInstance.Run()
	storage.Run()

	ctx = context.Background()
	lis = bufconn.Listen(bufSize)

	conn, err = grpc.DialContext(ctx, dialContext, grpc.WithDialer(BufDialer), grpc.WithInsecure())
	if !errors.IsEmpty(err) {
		panic(err)
	}

	s = grpc.NewServer()

	orderClient = pb.NewOrderHandlerClient(conn)

	// Register services
	channelService.RegisterStorage(storage)
	channelService.RegisterP2p(p2pInstance)

	joinres, _ := channelService.Join(ctx, &pb.JoinRequest{Asset: asset1, CounterAsset: asset2})
	channel = joinres.GetJoinedChannel()
}

func removeAllOrders() {
	storage.DeleteAllWithPrefix(string(interfaces.OrderPrefix))
}

func BufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func TestOrderStorageKeyPrefixer(t *testing.T) {
	prefixedBytes := getOrderStorageKey([]byte(asset1))
	assert.Equal(t, string(prefixedBytes), string(interfaces.OrderPrefix)+string(asset1))
}

func TestOrderCreation(t *testing.T) {
	createNewServerInstance()
	orderService.RegisterStorage(storage)
	orderService.RegisterP2p(p2pInstance)
	defer p2pInstance.Close()
	defer storage.Close()
	defer conn.Close()
	removeAllOrders()

	testOrder := pb.CreateRequest{ChannelID: channel.GetId(), Asset: asset1, CounterAsset: asset2, Amount: testAmount, Price: testPrice}

	var lastOrder *pb.Order

	// Register order endpoints with the gRPC server
	pb.RegisterOrderHandlerServer(s, orderService)

	go func() {
		if err := s.Serve(lis); !errors.IsEmpty(err) {
			t.Logf("Server exited with error: %v", err)
		}
		defer s.Stop()
	}()

	resp, err := orderClient.Create(ctx, &testOrder)
	assert.NoError(t, err)
	t.Logf("Created Order: %s", resp)
	assert.NotNil(t, resp)

	lastOrder = resp.GetCreatedOrder()
	storedOrder, err := orderClient.GetOrder(ctx, &pb.OrderSpecificRequest{OrderID: lastOrder.GetId(), ChannelID: channel.GetId()})
	assert.NoError(t, err)

	assert.Equal(t, lastOrder, storedOrder)

	resp2, err := orderClient.Delete(ctx, &pb.OrderSpecificRequest{OrderID: lastOrder.GetId(), ChannelID: channel.GetId()})
	assert.NoError(t, err)
	assert.NotNil(t, resp2)
}

func TestOrderReceive(t *testing.T) {
	createNewServerInstance()
	orderService.RegisterStorage(storage)
	defer p2pInstance.Close()
	defer storage.Close()
	defer conn.Close()
	removeAllOrders()

	testOrder := pb.CreateRequest{ChannelID: channel.GetId(), Asset: asset1, CounterAsset: asset2, Amount: testAmount, Price: testPrice}

	// Register order endpoints with the gRPC server
	pb.RegisterOrderHandlerServer(s, orderService)

	go func() {
		if err := s.Serve(lis); !errors.IsEmpty(err) {
			t.Fatalf("Server exited with error: %v", err)
		}
		defer s.Stop()
	}()

	order, err := orderService.Create(ctx, &testOrder)
	marshaledOrder, err := proto.Marshal(order)

	err = orderService.Receive(marshaledOrder)
	assert.NoError(t, err)

	storedOrder, err := orderClient.GetOrder(ctx, &pb.OrderSpecificRequest{OrderID: order.GetCreatedOrder().GetId()})
	assert.NoError(t, err)
	assert.NotNil(t, storedOrder)
}

func TestOrderGetAll(t *testing.T) {
	createNewServerInstance()
	orderService.RegisterStorage(storage)
	defer p2pInstance.Close()
	defer storage.Close()
	defer conn.Close()
	removeAllOrders()

	testOrder := pb.CreateRequest{ChannelID: channel.GetId(), Asset: asset1, CounterAsset: asset2, Amount: testAmount, Price: testPrice}

	// Register order endpoints with the gRPC server
	pb.RegisterOrderHandlerServer(s, orderService)

	go func() {
		if err := s.Serve(lis); !errors.IsEmpty(err) {
			t.Fatalf("Server exited with error: %v", err)
		}
		defer s.Stop()
	}()

	const testIterations = int(4)
	for i := 0; i < testIterations; i++ {
		_, err := orderClient.Create(ctx, &testOrder)
		assert.True(t, errors.IsEmpty(err))
	}

	resp, err := orderClient.GetAllOrders(ctx, &pb.Empty{})
	assert.True(t, errors.IsEmpty(err))
	orders := resp.GetOrders()
	assert.Equal(t, len(orders), testIterations)
}

package api

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/eqlabs/sprawl/db"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	bufconn "google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener
var client OrderHandlerClient
var ctx context.Context

var storage *db.Storage

// Construct a correct test Order
var testOrder = CreateRequest{Asset: []byte("ETH"), CounterAsset: []byte("BTC"), Amount: 52617562718, Price: 0.1}
var lastOrder *Order

func init() {
	storage := &db.Storage{}
	storage.SetDbPath("/var/lib/sprawl/test")
	storage.Run()

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	// Create an OrderService that stores the endpoints
	service := &OrderService{}
	// Register the storage service with it
	service.RegisterStorage(storage)

	RegisterOrderHandlerServer(s, service)

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

	client = NewOrderHandlerClient(conn)
}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func TestOrderCreation(t *testing.T) {
	resp, err := client.Create(ctx, &testOrder)
	assert.Equal(t, nil, err)
	t.Log("Created Order: ", resp)
	assert.NotEqual(t, false, resp)

	lastOrder = resp.GetCreatedOrder()

	resp2, err := client.Delete(ctx, &OrderSpecificRequest{Id: lastOrder.GetId()})
	assert.Equal(t, nil, err)
	assert.NotEqual(t, false, resp2)
}

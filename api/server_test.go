package api

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	bufconn "google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

var client OrderHandlerClient
var ctx context.Context

// Construct a correct test Order
var testOrder = CreateRequest{Asset: []byte("ETH"), CounterAsset: []byte("BTC"), Amount: 52617562718, Price: 0.1}
var lastOrder *Order

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	RegisterOrderHandlerServer(s, &OrderService{})
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
}

func TestOrderDeletion(t *testing.T) {
	resp, err := client.Delete(ctx, &OrderSpecificRequest{Id: lastOrder.GetId()})
	assert.Equal(t, nil, err)
	assert.NotEqual(t, false, resp)
}

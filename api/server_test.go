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

// Construct a correct test Order
var testOrder = CreateRequest{Asset: []byte("ETH"), CounterAsset: []byte("BTC"), Amount: 52617562718, Price: 0.1}

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	RegisterOrderHandlerServer(s, &OrderService{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(string, time.Duration) (net.Conn, error) {
	return lis.Dial()
}

func TestGrpc(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithDialer(bufDialer), grpc.WithInsecure())
	assert.Nil(t, err)

	defer conn.Close()
	client := NewOrderHandlerClient(conn)
	resp, err := client.Create(ctx, &testOrder)

	assert.Nil(t, err)
	t.Log(resp)
	assert.NotEqual(t, false, resp)
}

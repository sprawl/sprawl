package service

import (
	"context"
	"strconv"
	"testing"

	"github.com/sprawl/sprawl/pb"
	"github.com/sprawl/sprawl/util"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const serverTestKey string = "serverTestKey"
const serverTestEntry string = "serverTestEntry"
const apiPort string = "1337"
const serverAddr string = "localhost:1337"

func TestServerCreation(t *testing.T) {
	p2pInstance.Run()
	storage.Run()
	defer storage.Close()
	defer p2pInstance.Close()

	server := NewServer(nil, storage, p2pInstance)
	assert.Equal(t, server.Logger, new(util.PlaceholderLogger))

	server := NewServer(log, storage, p2pInstance, nil)
  
	assert.NotNil(t, server)
	assert.Equal(t, server.Logger, log)
	assert.Equal(t, server.Orders.Logger, log)
	assert.Equal(t, server.Orders.Storage, storage)
	assert.Equal(t, server.Channels.Storage, storage)
	assert.Equal(t, server.Orders.P2p, p2pInstance)
	assert.Equal(t, server.Channels.P2p, p2pInstance)

	var err error

	err = server.Orders.Storage.Put([]byte(serverTestKey), []byte(serverTestEntry))
	assert.NoError(t, err)
	server.Orders.Storage.DeleteAll()

	err = server.Channels.Storage.Put([]byte(serverTestKey), []byte(serverTestEntry))
	assert.NoError(t, err)
	server.Channels.Storage.DeleteAll()
}
func TestServerRun(t *testing.T) {
	p2pInstance.Run()
	storage.Run()
	defer storage.Close()
	defer p2pInstance.Close()

	server := NewServer(log, storage, p2pInstance, nil)
	port, err := strconv.ParseUint(apiPort, 10, 64)
	assert.NoError(t, err)
	go server.Run(uint(port))
	defer server.Close()

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewOrderHandlerClient(conn)
	resp, err := client.GetAllOrders(context.Background(), &pb.Empty{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

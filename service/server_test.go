package service

import (
	"context"
	"testing"

	"github.com/sprawl/sprawl/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const serverTestKey string = "serverTestKey"
const serverTestEntry string = "serverTestEntry"
const apiPort uint = 1337
const serverAddr string = "localhost:1337"

func TestServerCreation(t *testing.T) {
	p2pInstance.Run()
	storage.Run()
	defer storage.Close()
	defer p2pInstance.Close()

	server := NewServer(log, storage, p2pInstance)
	assert.NotNil(t, server)

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

	server := NewServer(log, storage, p2pInstance)
	go server.Run(apiPort)

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewOrderHandlerClient(conn)
	resp, err := client.GetAllOrders(context.Background(), &pb.Empty{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

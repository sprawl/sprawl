package service

import (
	"context"
	"testing"

	"github.com/eqlabs/sprawl/pb"
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

	server := NewServer(storage, p2pInstance)
	assert.NotEqual(t, server, nil)

	var err error

	err = server.Orders.Storage.Put([]byte(serverTestKey), []byte(serverTestEntry))
	assert.Equal(t, err, nil)
	server.Orders.Storage.DeleteAll()

	err = server.Channels.Storage.Put([]byte(serverTestKey), []byte(serverTestEntry))
	assert.Equal(t, err, nil)
	server.Channels.Storage.DeleteAll()
}

func TestServerRun(t *testing.T) {
	p2pInstance.Run()
	storage.Run()
	defer storage.Close()
	defer p2pInstance.Close()

	server := NewServer(storage, p2pInstance)
	go server.Run(apiPort)

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		t.Log(err)
	}
	defer conn.Close()

	client := pb.NewOrderHandlerClient(conn)
	resp, err := client.GetAllOrders(context.Background(), &pb.Empty{})
	if err != nil {
		t.Log(err)
	}
	assert.NotEqual(t, resp, nil)
}

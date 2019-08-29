package p2p

import (
	"context"
	"testing"

	"github.com/eqlabs/sprawl/pb"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/stretchr/testify/assert"
)

var test_channel pb.Channel = pb.Channel{Id: []byte("test_channel")}
var test_data []byte = []byte("test_data")

func TestCreateChannelString(t *testing.T) {
	assert.Equal(t, createChannelString(test_channel), string(test_channel.Id))
}

func TestInitContext(t *testing.T) {
	p2pInstance := NewP2p()
	p2pInstance.initContext()
	assert.Equal(t, p2pInstance.ctx, context.Background())
}

func TestInput(t *testing.T) {
	p2pInstance := NewP2p()
	go func() {
		p2pInstance.Input(test_data, test_channel)
	}()
	select {
	case message := <-p2pInstance.input:
		assert.Equal(t, *message.Channel, pb.Channel(test_channel))
		assert.Equal(t, message.Data, test_data)
	}
}

func TestPublish(t *testing.T) {
	p2pInstance := NewP2p()
	p2pInstance.initContext()
	p2pInstance.host, _ = libp2p.New(p2pInstance.ctx)
	p2pInstance.initPubSub()
	sub, _ := p2pInstance.ps.Subscribe(createChannelString(pb.Channel(test_channel)))
	go func() {
		p2pInstance.Input(test_data, test_channel)
	}()
	select {
	case message := <-p2pInstance.input:
		p2pInstance.handleInput(message)
		msg, _ := sub.Next(p2pInstance.ctx)
		assert.Equal(t, msg.Data, test_data)
	}
}

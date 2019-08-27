package p2p

import (
	"context"
	"testing"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/stretchr/testify/assert"
)

func TestCreateTopicString(t *testing.T) {
	assert.Equal(t, createTopicString("test_topic"), "/sprawl/test_topic")
}

func TestInitContext(t *testing.T) {
	p2pInstance := NewP2p()
	p2pInstance.initContext()
	assert.Equal(t, p2pInstance.ctx, context.Background())
}

func TestInput(t *testing.T) {
	p2pInstance := NewP2p()
	go func() {
		p2pInstance.Input([]byte("test_data"), "test_topic")
	}()
	select {
	case message := <-p2pInstance.input:
		assert.Equal(t, message.topic, "test_topic")
		assert.Equal(t, message.data, []byte("test_data"))
	}
}

func TestPublish(t *testing.T) {
	p2pInstance := NewP2p()
	p2pInstance.initContext()
	p2pInstance.host, _ = libp2p.New(p2pInstance.ctx)
	p2pInstance.initPubSub()
	sub, _ := p2pInstance.ps.Subscribe(createTopicString("test_topic"))
	go func() {
		p2pInstance.Input([]byte("test_data"), "test_topic")
	}()
	select {
	case message := <-p2pInstance.input:
		p2pInstance.handleInput(message)
		msg, _ := sub.Next(p2pInstance.ctx)
		assert.Equal(t, msg.Data, []byte("test_data"))
	}
}

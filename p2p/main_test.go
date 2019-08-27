package p2p

import (
	"context"
	"testing"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/stretchr/testify/assert"
)

const test_topic string = "test_topic"

var test_data []byte = []byte("test_data")

func TestCreateTopicString(t *testing.T) {
	assert.Equal(t, createTopicString(test_topic), baseTopic+test_topic)
}

func TestInitContext(t *testing.T) {
	p2pInstance := NewP2p()
	p2pInstance.initContext()
	assert.Equal(t, p2pInstance.ctx, context.Background())
}

func TestInput(t *testing.T) {
	p2pInstance := NewP2p()
	go func() {
		p2pInstance.Input(test_data, test_topic)
	}()
	select {
	case message := <-p2pInstance.input:
		assert.Equal(t, message.topic, test_topic)
		assert.Equal(t, message.data, test_data)
	}
}

func TestPublish(t *testing.T) {
	p2pInstance := NewP2p()
	p2pInstance.initContext()
	p2pInstance.host, _ = libp2p.New(p2pInstance.ctx)
	p2pInstance.initPubSub()
	sub, _ := p2pInstance.ps.Subscribe(createTopicString(test_topic))
	go func() {
		p2pInstance.Input(test_data, test_topic)
	}()
	select {
	case message := <-p2pInstance.input:
		p2pInstance.handleInput(message)
		msg, _ := sub.Next(p2pInstance.ctx)
		assert.Equal(t, msg.Data, test_data)
	}
}

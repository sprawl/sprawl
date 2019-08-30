package p2p

import (
	"context"
	"testing"

	"github.com/eqlabs/sprawl/pb"
	"github.com/gogo/protobuf/proto"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/stretchr/testify/assert"
)

var testChannel *pb.Channel = &pb.Channel{Id: []byte("testChannel")}
var testOrder *pb.Order = &pb.Order{Asset: string("ETH"), CounterAsset: string("BTC"), Amount: 52152, Price: 0.2, Id: []byte("jgkahgkjal")}
var testOrderInBytes []byte
var testWireMessage *pb.WireMessage

func TestCreateChannelString(t *testing.T) {
	assert.Equal(t, createChannelString(testChannel), string(testChannel.Id))
}

func TestInitContext(t *testing.T) {
	p2pInstance := NewP2p()
	p2pInstance.initContext()
	assert.Equal(t, p2pInstance.ctx, context.Background())
}

func TestSend(t *testing.T) {
	p2pInstance := NewP2p()
	testOrderInBytes, err := proto.Marshal(testOrder)
	if err != nil {
		panic(err)
	}
	testWireMessage = &pb.WireMessage{Channel: testChannel, Operation: pb.Operation_CREATE, Data: testOrderInBytes}
	p2pInstance.Send(testWireMessage)
	select {
	case message := <-p2pInstance.input:
		assert.Equal(t, *message.Channel, *testChannel)
	}
}

func TestPublish(t *testing.T) {
	p2pInstance := NewP2p()
	p2pInstance.initContext()
	p2pInstance.host, _ = libp2p.New(p2pInstance.ctx)
	p2pInstance.initPubSub()
	sub, _ := p2pInstance.ps.Subscribe(createChannelString(testChannel))
	testOrderInBytes, err := proto.Marshal(testOrder)
	if err != nil {
		panic(err)
	}
	testWireMessage = &pb.WireMessage{Channel: testChannel, Operation: pb.Operation_CREATE, Data: testOrderInBytes}
	p2pInstance.Send(testWireMessage)
	wireMessageAsBytes, err := proto.Marshal(testWireMessage)
	if err != nil {
		panic(err)
	}
	select {
	case message := <-p2pInstance.input:
		p2pInstance.handleInput(&message)
		msg, _ := sub.Next(p2pInstance.ctx)
		assert.Equal(t, msg.GetData(), wireMessageAsBytes)
	}
}

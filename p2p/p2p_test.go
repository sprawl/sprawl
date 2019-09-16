package p2p

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/identity"
	"github.com/eqlabs/sprawl/pb"
	"github.com/gogo/protobuf/proto"
	libp2p "github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const testConfigPath string = "../config/test"
const dbPathVar string = "database.path"

var testChannel *pb.Channel = &pb.Channel{Id: []byte("testChannel")}
var testOrder *pb.Order = &pb.Order{Asset: string("ETH"), CounterAsset: string("BTC"), Amount: 52152, Price: 0.2, Id: []byte("jgkahgkjal")}
var testOrderInBytes []byte
var testWireMessage *pb.WireMessage
var logger *zap.Logger
var log *zap.SugaredLogger
var testConfig *config.Config

func init() {
	logger, _ = zap.NewProduction()
	log = logger.Sugar()
	testConfig = &config.Config{Log: log}
}

func TestInitContext(t *testing.T) {
	privateKey, publicKey, err := identity.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err)
	p2pInstance := NewP2p(log, privateKey, publicKey)
	p2pInstance.initContext()
	assert.Equal(t, p2pInstance.ctx, context.Background())
}

func TestBootstrapping(t *testing.T) {
	privateKey, publicKey, err := identity.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err)
	p2pInstance := NewP2p(log, privateKey, publicKey)
	p2pInstance.addDefaultBootstrapPeers()
	var defaultBootstrapPeers addrList = dht.DefaultBootstrapPeers
	assert.Equal(t, p2pInstance.bootstrapPeers, defaultBootstrapPeers)
}

func TestSend(t *testing.T) {
	privateKey, publicKey, err := identity.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err)
	p2pInstance := NewP2p(log, privateKey, publicKey)

	testOrderInBytes, err := proto.Marshal(testOrder)
	assert.NoError(t, err)
	testWireMessage = &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}
	p2pInstance.Send(testWireMessage)

	message := <-p2pInstance.input
	assert.Equal(t, message.ChannelID, testChannel.GetId())
	assert.Equal(t, message.GetData(), testOrderInBytes)
}

func TestSubscription(t *testing.T) {
	privateKey, publicKey, err := identity.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err)
	p2pInstance := NewP2p(log, privateKey, publicKey)

	p2pInstance.initContext()
	p2pInstance.host, _ = libp2p.New(p2pInstance.ctx)

	p2pInstance.initPubSub()
	p2pInstance.Subscribe(testChannel)

	_, ok := p2pInstance.subscriptions[string(testChannel.GetId())]
	assert.True(t, ok)

	testOrderInBytes, err := proto.Marshal(testOrder)
	assert.NoError(t, err)
	testWireMessage = &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}

	go p2pInstance.Unsubscribe(testChannel)

	go func() {
		p2pInstance.Send(testWireMessage)
	}()

	go func() {
		p2pInstance.inputCheckLoop()
	}()
	<-p2pInstance.subscriptions[string(testChannel.GetId())]
}

func TestPublish(t *testing.T) {
	privateKey, publicKey, err := identity.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err)
	p2pInstance := NewP2p(log, privateKey, publicKey)

	p2pInstance.initContext()
	p2pInstance.host, _ = libp2p.New(p2pInstance.ctx)
	p2pInstance.initPubSub()

	sub, _ := p2pInstance.ps.Subscribe(string(testChannel.GetId()))
	testOrderInBytes, err := proto.Marshal(testOrder)
	assert.NoError(t, err)
	testWireMessage = &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}
	p2pInstance.Send(testWireMessage)
	wireMessageAsBytes, err := proto.Marshal(testWireMessage)
	assert.NoError(t, err)
	select {
	case message := <-p2pInstance.input:
		p2pInstance.handleInput(&message)
		msg, _ := sub.Next(p2pInstance.ctx)
		assert.Equal(t, msg.GetData(), wireMessageAsBytes)
	}
}

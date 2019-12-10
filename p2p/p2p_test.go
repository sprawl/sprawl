package p2p

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/golang/protobuf/proto"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/sprawl/sprawl/config"
	"github.com/sprawl/sprawl/identity"
	"github.com/sprawl/sprawl/pb"
	"github.com/sprawl/sprawl/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var testChannel *pb.Channel = &pb.Channel{Id: []byte("testChannel")}
var testOrder *pb.Order = &pb.Order{Asset: string("ETH"), CounterAsset: string("BTC"), Amount: 52152, Price: 0.2, Id: []byte("jgkahgkjal")}
var testOrderInBytes []byte
var testWireMessage *pb.WireMessage
var logger *zap.Logger
var log *zap.SugaredLogger
var testConfig *config.Config
var privateKey crypto.PrivKey
var publicKey crypto.PubKey

func init() {
	logger = zap.NewNop()
	log = logger.Sugar()
	testConfig = &config.Config{Logger: log}
	privateKey, publicKey, _ = identity.GenerateKeyPair(rand.Reader)
}

func TestServiceRegistration(t *testing.T) {
	p2pInstance := NewP2p(log, testConfig, privateKey, publicKey)
	orderService := &service.OrderService{}
	p2pInstance.AddReceiver(orderService)
	assert.Equal(t, orderService, p2pInstance.Receiver)
}

func TestInitContext(t *testing.T) {
	p2pInstance := NewP2p(log, testConfig, privateKey, publicKey)
	p2pInstance.initContext()
	assert.Equal(t, p2pInstance.ctx, context.Background())
}

func TestInitDHT(t *testing.T) {
	p2pInstance := NewP2p(log, testConfig, privateKey, publicKey)
	routing := p2pInstance.initDHT()
	assert.NotNil(t, routing)
}

func TestSend(t *testing.T) {
	p2pInstance := NewP2p(log, testConfig, privateKey, publicKey)

	testOrderInBytes, err := proto.Marshal(testOrder)
	assert.NoError(t, err)
	testWireMessage = &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}
	p2pInstance.Send(testWireMessage)

	message := <-p2pInstance.input
	assert.Equal(t, message.ChannelID, testChannel.GetId())
	assert.Equal(t, message.GetData(), testOrderInBytes)
}

func TestSubscription(t *testing.T) {
	p2pInstance := NewP2p(log, testConfig, privateKey, publicKey)

	p2pInstance.initContext()
	p2pInstance.host, _ = libp2p.New(p2pInstance.ctx)

	assert.Panics(t, func() { p2pInstance.Subscribe(testChannel) })

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
		p2pInstance.listenForInput()
	}()
	<-p2pInstance.subscriptions[string(testChannel.GetId())]
}

func TestPublish(t *testing.T) {
	p2pInstance := NewP2p(log, testConfig, privateKey, publicKey)

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

func TestRun(t *testing.T) {
	testConfig.ReadConfig(testConfigPath)
	p2pInstance := NewP2p(log, testConfig, privateKey, publicKey)
	// TODO: Acculi test this
	assert.NotPanics(t, p2pInstance.Run, "p2p run should not panic")
	assert.NotPanics(t, p2pInstance.Close, "p2p close should not panic")
}

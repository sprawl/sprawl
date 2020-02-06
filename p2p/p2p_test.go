package p2p

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/sprawl/sprawl/config"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/identity"
	"github.com/sprawl/sprawl/pb"
	"github.com/sprawl/sprawl/service"
	"github.com/sprawl/sprawl/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
var privateKey2 crypto.PrivKey
var publicKey2 crypto.PubKey

func init() {
	logger = zap.NewNop()
	log = logger.Sugar()
	testConfig = &config.Config{}
	testConfig.ReadConfig(testConfigPath)
	privateKey, publicKey, _ = identity.GenerateKeyPair(rand.Reader)
	privateKey2, publicKey2, _ = identity.GenerateKeyPair(rand.Reader)
}

type TestReceiver struct {
	t testing.T
	mock.Mock
}

func (r *TestReceiver) Receive(data []byte, from peer.ID) error {
	r.Called(data)
	return nil
}

func TestConstructor(t *testing.T) {
	orderService := &service.OrderService{}
	p2pInstance := NewP2p(testConfig, privateKey, publicKey, Logger(log), Receiver(orderService))
	assert.Equal(t, orderService, p2pInstance.Receiver)
	assert.Equal(t, log, p2pInstance.Logger)
	assert.Equal(t, p2pInstance.ctx, context.Background())
	p2pInstance = NewP2p(testConfig, privateKey, publicKey)
	assert.Equal(t, p2pInstance.Logger, &util.PlaceholderLogger{})
	assert.Nil(t, p2pInstance.Receiver)
	p2pInstance = NewP2p(testConfig, privateKey, publicKey, Receiver(nil))
	assert.Nil(t, p2pInstance.Receiver)
}

func TestInitDHT(t *testing.T) {
	p2pInstance := NewP2p(testConfig, privateKey, publicKey, Logger(log))
	routing := p2pInstance.initDHT()
	assert.NotNil(t, routing)
}

func TestSend(t *testing.T) {
	p2pInstance := NewP2p(testConfig, privateKey, publicKey, Logger(log))
	p2pInstance.InitHost(p2pInstance.CreateOptions()...)

	testOrderInBytes, err := proto.Marshal(testOrder)
	assert.NoError(t, err)

	testWireMessage = &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}
	p2pInstance.Send(testWireMessage)

	message := <-p2pInstance.input
	assert.Equal(t, message.ChannelID, testChannel.GetId())
	assert.Equal(t, message.GetData(), testOrderInBytes)
}

func TestSubscription(t *testing.T) {
	p2pInstance := NewP2p(testConfig, privateKey, publicKey, Logger(log))

	p2pInstance.host, _ = libp2p.New(p2pInstance.ctx)

	assert.Panics(t, func() { p2pInstance.Subscribe(testChannel) })

	p2pInstance.initPubSub()
	subCtx, _ := p2pInstance.Subscribe(testChannel)

	_, ok := p2pInstance.subscriptions[string(testChannel.GetId())]
	assert.True(t, ok)

	testOrderInBytes, err := proto.Marshal(testOrder)
	assert.NoError(t, err)

	testWireMessage = &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}

	p2pInstance.listenForInput()
	p2pInstance.Send(testWireMessage)

	p2pInstance.Unsubscribe(testChannel)

	<-subCtx.Done()
}

func TestPublish(t *testing.T) {
	p2pInstance := NewP2p(testConfig, privateKey, publicKey, Logger(log))

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
	p2pInstance := NewP2p(testConfig, privateKey, publicKey, Logger(log))
	assert.NotPanics(t, p2pInstance.Run, "p2p run should not panic")
	assert.NotPanics(t, p2pInstance.Close, "p2p close should not panic")
}

func TestChannelListener(t *testing.T) {
	// Initialize p2p instances
	p2pInstance1 := NewP2p(testConfig, privateKey, publicKey, Logger(log))
	p2pInstance2 := NewP2p(testConfig, privateKey2, publicKey2, Logger(log))

	testWireMessage = &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}
	wireMessageAsBytes, _ := proto.Marshal(testWireMessage)

	receiver := new(TestReceiver)
	receiver.Test(t)
	receiver.On("Receive", wireMessageAsBytes).Return(nil)
	p2pInstance2.AddReceiver(receiver)

	p2pInstance1.InitHost(p2pInstance1.CreateOptions()...)
	p2pInstance2.InitHost(p2pInstance2.CreateOptions()...)
	p2pInstance1.initPubSub()
	p2pInstance2.initPubSub()

	// Connect instances with each other
	err := p2pInstance1.host.Connect(p2pInstance1.ctx, p2pInstance2.GetAddrInfo())
	assert.NoError(t, err)
	err = p2pInstance2.host.Connect(p2pInstance2.ctx, p2pInstance1.GetAddrInfo())
	assert.NoError(t, err)

	subCtx1, err := p2pInstance1.Subscribe(testChannel)
	assert.NoError(t, err)
	subCtx2, err := p2pInstance2.Subscribe(testChannel)
	assert.NoError(t, err)

	p2pInstance1.Send(testWireMessage)

	p2pInstance1.Unsubscribe(testChannel)
	p2pInstance2.Unsubscribe(testChannel)

	<-subCtx1.Done()
	<-subCtx2.Done()
}

func TestStreams(t *testing.T) {
	// Initialize p2p instances
	p2pInstance1 := NewP2p(testConfig, privateKey, publicKey, Logger(log))
	p2pInstance2 := NewP2p(testConfig, privateKey2, publicKey2, Logger(log))

	testWireMessage = &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}
	wireMessageAsBytes, _ := proto.Marshal(testWireMessage)
	receiver := new(TestReceiver)
	receiver.Test(t)
	receiver.On("Receive", wireMessageAsBytes).Return(nil)
	p2pInstance2.AddReceiver(receiver)

	p2pInstance1.InitHost(p2pInstance1.CreateOptions()...)
	p2pInstance2.InitHost(p2pInstance2.CreateOptions()...)

	// Connect instances with each other
	err := p2pInstance1.host.Connect(p2pInstance1.ctx, p2pInstance2.GetAddrInfo())
	assert.NoError(t, err)
	err = p2pInstance2.host.Connect(p2pInstance2.ctx, p2pInstance1.GetAddrInfo())
	assert.NoError(t, err)

	peerList := p2pInstance1.GetAllPeers()
	assert.NotEmpty(t, peerList)
	assert.Contains(t, peerList, p2pInstance2.GetHostID())

	// Open bilateral stream
	stream, _ := p2pInstance1.OpenStream(p2pInstance2.GetHostID())

	// Assert p2p.streams lengths
	assert.Len(t, p2pInstance1.streams, 1)
	assert.Len(t, p2pInstance2.streams, 0)
	keys := []string{}
	for key := range p2pInstance1.streams {
		keys = append(keys, key)
	}
	assert.Equal(t, keys[0], p2pInstance2.GetHostIDString())

	// Write from p2pInstance1 to p2pInstance2
	err = stream.WriteToStream(wireMessageAsBytes)
	time.Sleep(time.Second / 2)
	assert.True(t, errors.IsEmpty(err))

	// Check that the message was received on p2pInstance2's end
	receiver.AssertCalled(t, "Receive", wireMessageAsBytes)

	// Close the stream on p2pInstance1's end
	p2pInstance1.CloseStream(p2pInstance2.GetHostID())
	assert.Len(t, p2pInstance1.streams, 0)
}

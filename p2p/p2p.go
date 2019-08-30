package p2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/eqlabs/sprawl/interfaces"

	"github.com/eqlabs/sprawl/pb"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	routing "github.com/libp2p/go-libp2p-core/routing"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	config "github.com/libp2p/go-libp2p/config"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// A new type we need for writing a custom flag parser
type addrList []multiaddr.Multiaddr

const baseTopic = "/sprawl/"

// P2p stores all things required to converse with other peers in the Sprawl network and save data locally
type P2p struct {
	ps               *pubsub.PubSub
	ctx              context.Context
	host             host.Host
	kademliaDHT      *dht.IpfsDHT
	routingDiscovery *discovery.RoutingDiscovery
	peerChan         <-chan peer.AddrInfo
	bootstrapPeers   addrList
	input            chan pb.WireMessage
	orders           interfaces.OrderService
	channels         interfaces.ChannelService
}

func NewP2p() (p2p *P2p) {
	p2p = &P2p{
		input: make(chan pb.WireMessage),
	}
	return
}

func (p2p *P2p) InputCheckLoop() (err error) {
	for {
		select {
		case message := <-p2p.input:
			p2p.handleInput(message)
		}
	}
}

// RegisterOrderService registers an order service to persist order data locally
func (p2p *P2p) RegisterOrderService(orders interfaces.OrderService) {
	p2p.orders = orders
}

// RegisterChannelService registers a channel service to persist joined channels locally
func (p2p *P2p) RegisterChannelService(channels interfaces.ChannelService) {
	p2p.channels = channels
}

func (p2p *P2p) handleInput(message pb.WireMessage) {
	err := p2p.ps.Publish(createChannelString(message.Channel), message.Data)
	if err != nil {
		fmt.Printf("Error publishing with %s, %v", message.Data, err)
	}
}

// Send queues a message for sending to other peers
func (p2p *P2p) Send(channel *pb.Channel, data []byte) {
	go func() {
		p2p.input <- pb.WireMessage{Channel: channel, Data: data}
	}()
}

func createChannelString(channel *pb.Channel) string {
	return string(channel.GetId())
}

func (p2p *P2p) initPubSub() {
	var err error
	p2p.ps, err = pubsub.NewGossipSub(p2p.ctx, p2p.host)
	if err != nil {
		panic(err)
	}
}

// Subscribe subscribes to a libp2p pubsub channel defined with "channel"
func (p2p *P2p) Subscribe(channel *pb.Channel) {
	sub, err := p2p.ps.Subscribe(createChannelString(channel))
	if err != nil {
		panic(err)
	}
	go func(ctx context.Context) {
		for {
			msg, err := sub.Next(p2p.ctx)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Message received: %s\n", msg)
			p2p.orders.Receive(msg.GetData())
		}
	}(p2p.ctx)
}

func (p2p *P2p) initContext() {
	p2p.ctx = context.Background()
}

func (p2p *P2p) bootstrapDHT() {
	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	var err error
	if err = p2p.kademliaDHT.Bootstrap(p2p.ctx); err != nil {
		panic(err)
	}
}

func (p2p *P2p) initBootstrapPeers(bootstrapPeers addrList) {
	p2p.bootstrapPeers = bootstrapPeers
}

func (p2p *P2p) addDefaultBootstrapPeers() {
	p2p.initBootstrapPeers(dht.DefaultBootstrapPeers)
}

func (p2p *P2p) connectToPeers() {
	var wg sync.WaitGroup
	for _, peerAddr := range p2p.bootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := p2p.host.Connect(p2p.ctx, *peerinfo); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Connected to : %s\n", peerinfo)
			}
		}()
	}
	wg.Wait()
}

func (p2p *P2p) createRoutingDiscovery() {
	p2p.routingDiscovery = discovery.NewRoutingDiscovery(p2p.kademliaDHT)
}

func (p2p *P2p) advertise() {
	discovery.Advertise(p2p.ctx, p2p.routingDiscovery, baseTopic)
}

func (p2p *P2p) findPeers() {
	var err error
	p2p.peerChan, err = p2p.routingDiscovery.FindPeers(p2p.ctx, baseTopic)
	if err != nil {
		panic(err)
	}
}

func (p2p *P2p) initDHT() config.Option {
	NewDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		p2p.kademliaDHT, err = dht.New(p2p.ctx, h)
		return p2p.kademliaDHT, err
	}
	return libp2p.Routing(NewDHT)

}

func (p2p *P2p) initHost(routing config.Option) {
	var err error
	p2p.host, err = libp2p.New(p2p.ctx,
		routing,
		libp2p.EnableRelay(),
		libp2p.EnableAutoRelay(),
		libp2p.NATPortMap(),
	)
	if err != nil {
		panic(err)
	}
}

// Run runs the p2p network
func (p2p *P2p) Run() {
	p2p.initContext()
	p2p.initHost(p2p.initDHT())
	p2p.addDefaultBootstrapPeers()
	p2p.connectToPeers()
	p2p.createRoutingDiscovery()
	p2p.advertise()
	p2p.findPeers()
	p2p.initPubSub()
	p2p.bootstrapDHT()
}

// Close closes the underlying libp2p host
func (p2p *P2p) Close() {
	p2p.host.Close()
}

package p2p

import (
	"context"
	"fmt"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	routing "github.com/libp2p/go-libp2p-core/routing"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	configt "github.com/libp2p/go-libp2p/config"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// A new type we need for writing a custom flag parser
type addrList []multiaddr.Multiaddr
const BaseTopic = "/sprawl/"

func createTopicString(topic string) string {
	return BaseTopic + topic
}

func PublishMessage(ps *pubsub.PubSub, topic string, input []byte) {
	err := ps.Publish(createTopicString(topic), input)
	if err != nil {
		panic(err)
	}
}

func (p2p *P2p) initPubSub(ctx context.Context, host host.Host) {
	var err error
	p2p.ps, err = pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}
}

func Subscribe(ps *pubsub.PubSub, ctx context.Context, topic string) {
	sub, err := ps.Subscribe(createTopicString(topic))
	if err != nil {
		panic(err)
	}
	go func(ctx context.Context) {
		for {
			msg, err := sub.Next(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println("Message: %s", msg)
		}
	}(ctx)
}

type P2p struct {
	ps               *pubsub.PubSub
	ctx              context.Context
	host             host.Host
	kademliaDHT      *dht.IpfsDHT
	routingDiscovery *discovery.RoutingDiscovery
	peerChan         <-chan peer.AddrInfo
	bootstrapPeers   addrList
}

func (p2p *P2p) initContext() {
	p2p.ctx = context.Background()
}

func bootstrapDHT(ctx context.Context, kademliaDHT *dht.IpfsDHT) {
	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	var err error
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
}

func (p2p *P2p) initBootstrapPeers(bootstrapPeers addrList) {
	p2p.bootstrapPeers = bootstrapPeers
}

func (p2p *P2p) addDefaultBootstrapPeers() {
	p2p.initBootstrapPeers(dht.DefaultBootstrapPeers)
}

func connectToPeers(ctx context.Context, host host.Host, bootstrapPeers addrList) {
	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Connected to : %s", peerinfo)
			}
		}()
	}
	wg.Wait()
}

func (p2p *P2p) createRoutingDiscovery(kademliaDHT *dht.IpfsDHT) {
	p2p.routingDiscovery = discovery.NewRoutingDiscovery(kademliaDHT)
}

func advertise(ctx context.Context, routingDiscovery *discovery.RoutingDiscovery) {
	discovery.Advertise(ctx, routingDiscovery, BaseTopic)
}

func (p2p *P2p) findPeers(ctx context.Context, routingDiscovery *discovery.RoutingDiscovery) {
	var err error
	p2p.peerChan, err = routingDiscovery.FindPeers(ctx, BaseTopic)
	if err != nil {
		panic(err)
	}
}

func (p2p *P2p) initDHT() configt.Option {
	NewDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		p2p.kademliaDHT, err = dht.New(p2p.ctx, h)
		return p2p.kademliaDHT, err
	}
	return libp2p.Routing(NewDHT)

}

func (p2p *P2p) initHost(routing configt.Option) {
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
	connectToPeers(p2p.ctx, p2p.host, p2p.bootstrapPeers)
	p2p.createRoutingDiscovery(p2p.kademliaDHT)
	advertise(p2p.ctx, p2p.routingDiscovery)
	p2p.findPeers(p2p.ctx, p2p.routingDiscovery)
	p2p.initPubSub(p2p.ctx, p2p.host)
	bootstrapDHT(p2p.ctx, p2p.kademliaDHT)
	select {}
}

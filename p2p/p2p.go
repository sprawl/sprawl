package p2p

import (
	"context"
	"sync"
	"time"

	"github.com/eqlabs/sprawl/interfaces"
	"github.com/gogo/protobuf/proto"

	"github.com/eqlabs/sprawl/pb"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
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
	Logger           interfaces.Logger
	privateKey       crypto.PrivKey
	publicKey        crypto.PubKey
	ps               *pubsub.PubSub
	ctx              context.Context
	host             host.Host
	kademliaDHT      *dht.IpfsDHT
	routingDiscovery *discovery.RoutingDiscovery
	peerChan         <-chan peer.AddrInfo
	bootstrapPeers   addrList
	input            chan pb.WireMessage
	subscriptions    map[string]chan bool
	Orders           interfaces.OrderService
	Channels         interfaces.ChannelService
}

// NewP2p returns a P2p struct with an input channel
func NewP2p(log interfaces.Logger, privateKey crypto.PrivKey, publicKey crypto.PubKey) (p2p *P2p) {
	p2p = &P2p{
		Logger:        log,
		privateKey:    privateKey,
		publicKey:     publicKey,
		input:         make(chan pb.WireMessage),
		subscriptions: make(map[string]chan bool),
	}
	return
}

func (p2p *P2p) inputCheckLoop() (err error) {
	for {
		select {
		case message := <-p2p.input:
			p2p.handleInput(&message)
		}
	}
}

func (p2p *P2p) checkForPeers() {
	if p2p.Logger != nil {
		p2p.Logger.Infof("This node's ID: %s\n", p2p.host.ID())
	}
	go func(ctx context.Context) {
		for peer := range p2p.peerChan {
			if peer.ID == p2p.host.ID() {
				if p2p.Logger != nil {
					p2p.Logger.Debug("Found a new peer!")
					p2p.Logger.Debug("But the peer was you!")
				}
				continue
			}
			if p2p.Logger != nil {
				p2p.Logger.Infof("Found a new peer: %s\n", peer.ID)
			}
			p2p.ps.ListPeers(baseTopic)
			if err := p2p.host.Connect(ctx, peer); err != nil {
				if p2p.Logger != nil {
					p2p.Logger.Error(err)
				}
			} else {
				if p2p.Logger != nil {
					p2p.Logger.Infof("Connected to: %s\n", peer)
				}
			}
		}
	}(p2p.ctx)
}

// RegisterOrderService registers an order service to persist order data locally
func (p2p *P2p) RegisterOrderService(orders interfaces.OrderService) {
	p2p.Orders = orders
}

// RegisterChannelService registers a channel service to persist joined channels locally
func (p2p *P2p) RegisterChannelService(channels interfaces.ChannelService) {
	p2p.Channels = channels
}

func (p2p *P2p) handleInput(message *pb.WireMessage) {
	buf, err := proto.Marshal(message)
	err = p2p.ps.Publish(string(message.GetChannelID()), buf)
	if err != nil {
		if p2p.Logger != nil {
			p2p.Logger.Errorf("Error publishing with %s, %v", message.Data, err)
		}
	}
}

// Send queues a message for sending to other peers
func (p2p *P2p) Send(message *pb.WireMessage) {
	if p2p.Logger != nil {
		p2p.Logger.Debugf("Sending order %s to channel %s", message.GetData(), message.GetChannelID())
	}
	go func(ctx context.Context) {
		p2p.input <- *message
	}(p2p.ctx)
}

func (p2p *P2p) initPubSub() {
	var err error
	p2p.ps, err = pubsub.NewGossipSub(p2p.ctx, p2p.host)
	if err != nil {
		if p2p.Logger != nil {
			p2p.Logger.Error(err)
		}
	}
}

// Subscribe subscribes to a libp2p pubsub channel defined with "channel"
func (p2p *P2p) Subscribe(channel *pb.Channel) {
	if p2p.Logger != nil {
		p2p.Logger.Infof("Subscribing to channel %s with options: %s", channel.GetId(), channel.GetOptions())
	}
	sub, err := p2p.ps.Subscribe(string(channel.GetId()))
	if err != nil {
		if p2p.Logger != nil {
			p2p.Logger.Error(err)
		}
	}

	quitSignal := make(chan bool)
	p2p.subscriptions[string(channel.GetId())] = quitSignal

	go func(ctx context.Context) {
		for {
			msg, err := sub.Next(ctx)
			if err != nil {
				if p2p.Logger != nil {
					p2p.Logger.Error(err)
				}
			}

			data := msg.GetData()
			peer := msg.GetFrom()

			if peer != p2p.host.ID() {
				if p2p.Logger != nil {
					p2p.Logger.Infof("Received order from peer %s: %s", peer, data)
				}

				if p2p.Orders != nil {
					err = p2p.Orders.Receive(data)
					if err != nil {
						if p2p.Logger != nil {
							p2p.Logger.Error(err)
						}
					}
				} else {
					if p2p.Logger != nil {
						p2p.Logger.Warn("P2p: OrderService not registered with p2p, not persisting incoming orders to DB!")
					}
				}
			}

			select {
			case quit := <-quitSignal: //Delete subscription
				if quit {
					delete(p2p.subscriptions, string(channel.GetId()))
					return
				}
			default:
			}
		}
	}(p2p.ctx)
}

// Unsubscribe sends a quit signal to a channel goroutine
func (p2p *P2p) Unsubscribe(channel *pb.Channel) {
	p2p.subscriptions[string(channel.GetId())] <- true
}

func (p2p *P2p) initContext() {
	p2p.ctx = context.Background()
}

func (p2p *P2p) bootstrapDHT() {
	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	var err error

	bootstrapConfig := dht.BootstrapConfig{
		Queries: 1,
		Period:  time.Duration(2 * time.Minute),
		Timeout: time.Duration(10 * time.Second),
	}

	err = p2p.kademliaDHT.BootstrapWithConfig(p2p.ctx, bootstrapConfig)

	if err != nil {
		if p2p.Logger != nil {
			p2p.Logger.Error(err)
		}
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
	if p2p.Logger != nil {
		p2p.Logger.Info("Connecting to bootstrap peers")
	}

	for _, peerAddr := range p2p.bootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)

		go func() {
			defer wg.Done()
			if err := p2p.host.Connect(p2p.ctx, *peerinfo); err != nil {
				if p2p.Logger != nil {
					p2p.Logger.Warnf("Error connecting to bootstrap peer %s", err)
				}
			} else {
				if p2p.Logger != nil {
					p2p.Logger.Infof("Connected to node: %s\n", peerinfo)
				}
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
		if p2p.Logger != nil {
			p2p.Logger.Error(err)
		}
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

func (p2p *P2p) initHost(options ...config.Option) {
	var err error
	p2p.host, err = libp2p.New(
		p2p.ctx,
		options...)
	if err != nil {
		if p2p.Logger != nil {
			p2p.Logger.Error(err)
		}
	}
}

// Run runs the p2p network
func (p2p *P2p) Run() {
	p2p.initContext()
	p2p.initHost(p2p.CreateOptions()...)//p2p.initDHT())
	p2p.addDefaultBootstrapPeers()
	p2p.connectToPeers()
	p2p.createRoutingDiscovery()
	p2p.advertise()
	p2p.findPeers()
	p2p.initPubSub()
	p2p.bootstrapDHT()
	go func() {
		p2p.inputCheckLoop()
	}()
	p2p.checkForPeers()
}

// Close closes the underlying libp2p host
func (p2p *P2p) Close() {
	p2p.host.Close()
}

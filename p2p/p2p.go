package p2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/sprawl/sprawl/interfaces"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2pConfig "github.com/libp2p/go-libp2p/config"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/pb"
)

const networkID = "/sprawl/"

// P2p stores all things required to converse with other peers in the Sprawl network and save data locally
type P2p struct {
	Config           interfaces.Config
	privateKey       crypto.PrivKey
	publicKey        crypto.PubKey
	ps               *pubsub.PubSub
	ctx              context.Context
	host             host.Host
	kademliaDHT      *dht.IpfsDHT
	routingDiscovery *discovery.RoutingDiscovery
	peerChan         <-chan peer.AddrInfo
	input            chan pb.WireMessage
	subscriptions    map[string]chan bool
	streams          map[string]*Stream
	Logger           interfaces.Logger
	storage          interfaces.Storage
	Receiver         interfaces.Receiver
}

// NewP2p returns a P2p struct with an input channel
func NewP2p(config interfaces.Config, privateKey crypto.PrivKey, publicKey crypto.PubKey, opts ...Option) (p2p *P2p) {
	p2p = &P2p{
		Config:        config,
		privateKey:    privateKey,
		publicKey:     publicKey,
		input:         make(chan pb.WireMessage),
		subscriptions: make(map[string]chan bool),
		streams:       make(map[string]*Stream),
	}

	for _, opt := range opts {
		err := opt(p2p)
		if err != nil {
			return nil
		}
	}

	return p2p
}

// AddReceiver registers a data receiver function with p2p
func (p2p *P2p) AddReceiver(receiver interfaces.Receiver) {
	p2p.Receiver = receiver
}

// InitContext ties background context to p2p.ctx
func (p2p *P2p) InitContext() {
	p2p.ctx = context.Background()
}

// InitHost creates a libp2p host with given options
func (p2p *P2p) InitHost(options ...libp2pConfig.Option) {
	var err error

	// Construct the libp2p host with options
	p2p.host, err = libp2p.New(
		p2p.ctx,
		options...)

	// Set stream handler for libp2p host
	p2p.host.SetStreamHandler(networkID, p2p.handleStream)

	if !errors.IsEmpty(err) {
		if p2p.Logger != nil {
			p2p.Logger.Error(errors.E(errors.Op("Creating host"), err))
		}
	}

	err = p2p.kademliaDHT.Bootstrap(p2p.ctx)

	if !errors.IsEmpty(err) {
		if p2p.Logger != nil {
			p2p.Logger.Error(errors.E(errors.Op("Constructing DHT"), err))
		}
	}
}

// GetHostIDString returns the underlying libp2p host's peer.ID as a string
func (p2p *P2p) GetHostIDString() string {
	return p2p.host.ID().String()
}

// GetHostID returns the underlying libp2p host's peer.ID
func (p2p *P2p) GetHostID() peer.ID {
	return p2p.host.ID()
}

// GetAddrInfo uses p2p.ConstructAddrInfo to get this peer's own AddrInfo
func (p2p *P2p) GetAddrInfo() peer.AddrInfo {
	return p2p.ConstructAddrInfo(p2p.GetHostID(), p2p.host.Addrs())
}

// ConstructAddrInfo is used to construct peer.AddrInfo especially in tests
func (p2p *P2p) ConstructAddrInfo(id peer.ID, addrs []ma.Multiaddr) peer.AddrInfo {
	return peer.AddrInfo{ID: id, Addrs: addrs}
}

func (p2p *P2p) initPubSub() {
	var err error
	p2p.ps, err = pubsub.NewGossipSub(p2p.ctx, p2p.host)
	if !errors.IsEmpty(err) {
		if p2p.Logger != nil {
			p2p.Logger.Error(err)
		}
	}
}

func (p2p *P2p) connectToNetwork() {
	var wg sync.WaitGroup
	if p2p.Logger != nil {
		p2p.Logger.Info("Connecting to bootstrap peers")
	}

	for _, peerAddr := range p2p.defaultBootstrapPeers() {
		// Parse URLs from each bootstrap peer
		peerinfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil && p2p.Logger != nil {
			p2p.Logger.Errorf("Bootstrap peer multiaddress %s is invalid: %s", peerAddr, err)
		} else {
			// Connect to the peer synchronically if the URL is correct
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := p2p.host.Connect(p2p.ctx, *peerinfo); !errors.IsEmpty(err) {
					if p2p.Logger != nil {
						p2p.Logger.Debugf("Error connecting to bootstrap peer %s", err)
					} else {
						p2p.Logger.Debugf("Successfully connected to bootstrap peer %s", peerinfo)
					}
				}
			}()
		}
	}

	wg.Wait()
}

func (p2p *P2p) startDiscovery() {
	// Add Kademlia routing discovery
	p2p.routingDiscovery = discovery.NewRoutingDiscovery(p2p.kademliaDHT)

	// Start the advertiser service
	discovery.Advertise(p2p.ctx, p2p.routingDiscovery, networkID)

	var err error
	// Ingest newly found peers into p2p.peerChan
	p2p.peerChan, err = p2p.routingDiscovery.FindPeers(p2p.ctx, networkID)

	if !errors.IsEmpty(err) {
		if p2p.Logger != nil {
			p2p.Logger.Error(errors.E(errors.Op("Find peers"), err))
		}
	}
}

func (p2p *P2p) listenForPeers() {
	if p2p.Logger != nil {
		p2p.Logger.Infof("This node's ID: %s\n", p2p.host.ID())
		p2p.Logger.Infof("Listening to the following addresses: %s\n", p2p.host.Addrs())
	}
	var wg sync.WaitGroup

	go func(ctx context.Context) {
		for peer := range p2p.peerChan {
			if peer.ID == p2p.host.ID() {
				if p2p.Logger != nil {
					p2p.Logger.Debug("Found yourself!")
				}
				continue
			}
			if p2p.Logger != nil {
				p2p.Logger.Infof("Found a new peer: %s\n", peer.ID)
			}

			// Waits on each peerInfo until they are connected or the connection failed
			wg.Add(1)
			go func(ctx context.Context) {
				defer wg.Done()
				if err := p2p.host.Connect(ctx, peer); !errors.IsEmpty(err) {
					if p2p.Logger != nil {
						p2p.Logger.Error(errors.E(errors.Op("Connect"), err))
					}
				} else {
					if p2p.Logger != nil {
						p2p.Logger.Infof("Connected to: %s\n", peer)
					}
				}
			}(p2p.ctx)
			wg.Wait()
		}
	}(p2p.ctx)
}

// handleInput takes in any local input, marshals it to Protobuf bytes and publishes it
func (p2p *P2p) handleInput(message *pb.WireMessage) {
	buf, err := proto.Marshal(message)
	if !errors.IsEmpty(err) {
		if p2p.Logger != nil {
			p2p.Logger.Error(errors.E(errors.Op("Marshal proto"), err))
		}
	}
	p2p.Logger.Debugf("Publishing to topic %s!", string(message.GetChannelID()))
	err = p2p.ps.Publish(string(message.GetChannelID()), buf)
	if !errors.IsEmpty(err) {
		if p2p.Logger != nil {
			p2p.Logger.Error(errors.E(errors.Op("Marshal proto"), fmt.Sprintf("%v, message data: %s", err.Error(), message.Data)))
		}
	}
}

// listenForInput pushes new items in channel p2p.input to p2p.handleInput
func (p2p *P2p) listenForInput() {
	go func() {
		for {
			select {
			case message := <-p2p.input:
				p2p.handleInput(&message)
			}
		}
	}()
}

// Send queues a message for sending to other peers
func (p2p *P2p) Send(message *pb.WireMessage) {
	go func(ctx context.Context) {
		p2p.input <- *message
	}(p2p.ctx)
}

func (p2p *P2p) GetAllPeers() []string {
	peers := p2p.host.Network().Peers()
	peersList := make([]string, len(peers))
	for _, value := range peersList {
		peersList = append(peersList, value)
	}
	return peersList
}

func (p2p *P2p) BlacklistPeer(pbPeer *pb.Peer) {
	peer, _ := peer.IDFromString(pbPeer.Id)
	p2p.ps.BlacklistPeer(peer)
}

// Subscribe subscribes to a libp2p pubsub channel defined with "channel"
func (p2p *P2p) Subscribe(channel *pb.Channel) {
	if p2p.Logger != nil {
		p2p.Logger.Infof("Subscribing to channel %s with options: %s", channel.GetId(), channel.GetOptions())
	}

	topic, err := p2p.ps.Join(string(channel.GetId()))
	if !errors.IsEmpty(err) {
		if p2p.Logger != nil {
			p2p.Logger.Error(errors.E(errors.Op("Join libp2p Topic"), err))
		}
	}
	sub, err := topic.Subscribe()
	if !errors.IsEmpty(err) {
		if p2p.Logger != nil {
			p2p.Logger.Error(errors.E(errors.Op("Subscribe to libp2p Topic"), err))
		}
	}

	quitSignal := make(chan bool)
	p2p.subscriptions[string(channel.GetId())] = quitSignal

	// Listen for new data
	p2p.listenToChannel(sub, channel, quitSignal)

	// Listen for new topic subscribers
	p2p.pingNewMembers(sub.Topic(), topic)
}

// Unsubscribe sends a quit signal to a channel goroutine
func (p2p *P2p) Unsubscribe(channel *pb.Channel) {
	p2p.subscriptions[string(channel.GetId())] <- true
}

// Run runs the p2p network
func (p2p *P2p) Run() {
	p2p.InitContext()

	// Initialize the p2p host with options
	p2p.InitHost(p2p.CreateOptions()...)

	// Connect to Sprawl & IPFS main nodes for peer discovery
	p2p.connectToNetwork()

	// Start finding peers on the network
	p2p.startDiscovery()

	// Start PubSub
	p2p.initPubSub()

	// Listen for local and network input
	p2p.listenForInput()

	// Continuously connect to other Sprawl peers
	p2p.listenForPeers()
}

// Close closes the underlying libp2p host
func (p2p *P2p) Close() {
	p2p.Logger.Debug("P2P shutting down")
	p2p.host.Close()
}

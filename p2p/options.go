package p2p

import (
	"fmt"
	"strconv"

	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/interfaces"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	routing "github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	libp2pConfig "github.com/libp2p/go-libp2p/config"
	ma "github.com/multiformats/go-multiaddr"
)

const addrTemplate string = "/ip4/%s/tcp/%s"

// Options for this p2p package, unrelated to libp2pConfig.Option
type Options struct {
	Logger  interfaces.Logger
	Storage interfaces.Storage
}

// Option type that allows us to have many underlying types of options.
type Option func(*P2p) error

// Storage is an interface to a data store for the p2p package
func Storage(storage interfaces.Storage) Option {
	return func(p *P2p) error {
		p.storage = storage
		return nil
	}
}

// Logger is an interface to a logger for the p2p package
func Logger(logger interfaces.Logger) Option {
	return func(p *P2p) error {
		p.Logger = logger
		return nil
	}
}

// Receiver receives all data that other peers send on pubsub channels
func Receiver(receiver interfaces.Receiver) Option {
	return func(p *P2p) error {
		p.Receiver = receiver
		return nil
	}
}

func (p2p *P2p) defaultBootstrapPeers() []ma.Multiaddr {
	peers := []ma.Multiaddr{}
	peers = append(peers, dht.DefaultBootstrapPeers...)
	/* sprawlBootstrapAddresses := []string{"/dnsaddr/bootstrap.sprawl.equilibrium.co"}
	for _, addr := range sprawlBootstrapAddresses {
		mAddr, _ := ma.NewMultiaddr(addr)
		peers = append(peers, mAddr)
	} */
	return peers
}

func createMultiAddr(externalIP string, p2pPort string) (ma.Multiaddr, error) {
	return ma.NewMultiaddr(fmt.Sprintf(addrTemplate, externalIP, p2pPort))
}

func (p2p *P2p) initDHT() libp2pConfig.Option {
	NewDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		p2p.kademliaDHT, err = dht.New(p2p.ctx, h)
		if !errors.IsEmpty(err) {
			p2p.Logger.Error(errors.E(errors.Op("Add dht"), err))
		}
		return p2p.kademliaDHT, err
	}
	return libp2p.Routing(NewDHT)
}

// CreateOptions queries p2p.Config for any user-submitted options and assigns defaults
func (p2p *P2p) CreateOptions() []libp2pConfig.Option {
	options := []libp2pConfig.Option{}
	externalIP := p2p.Config.GetExternalIP()
	p2pPort := p2p.Config.GetP2PPort()

	// Non-configurable options, since we always need an identity and the DHT discovery
	options = append(options, p2p.initDHT())
	options = append(options, libp2p.Identity(p2p.privateKey))

	// libp2p relay options
	if p2p.Config.GetRelaySetting() {
		options = append(options, libp2p.EnableRelay())
	}
	if p2p.Config.GetAutoRelaySetting() {
		options = append(options, libp2p.EnableAutoRelay())
	}

	// If NAT port map is not enabled, define listened addresses and port manually
	if p2p.Config.GetNATPortMapSetting() {
		options = append(options, libp2p.NATPortMap())
	} else {
		multiaddrs := []ma.Multiaddr{}
		if externalIP != "" {
			extMultiAddr, err := createMultiAddr(externalIP, strconv.FormatUint(uint64(p2pPort), 10))
			if !errors.IsEmpty(err) {
				p2p.Logger.Error(errors.E(errors.Op("Creating multiaddr"), err))
			}
			multiaddrs = append(multiaddrs, extMultiAddr)
		}
		addrFactory := func(addrs []ma.Multiaddr) []ma.Multiaddr {
			return multiaddrs
		}
		options = append(options, libp2p.ListenAddrs(multiaddrs...))
		options = append(options, libp2p.AddrsFactory(addrFactory))
	}

	return options
}

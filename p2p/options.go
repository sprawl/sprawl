package p2p

import (
	"fmt"

	"github.com/sprawl/sprawl/errors"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	routing "github.com/libp2p/go-libp2p-routing"
	libp2pConfig "github.com/libp2p/go-libp2p/config"
	ma "github.com/multiformats/go-multiaddr"
)

const addrTemplate string = "/ip4/%s/tcp/%s"

func defaultListenAddrs(p2pPort string) []ma.Multiaddr {
	multiaddrs := []ma.Multiaddr{}
	localhost, _ := ma.NewMultiaddr(fmt.Sprintf(addrTemplate, "0.0.0.0", p2pPort))
	multiaddrs = append(multiaddrs, localhost)
	return multiaddrs
}

func defaultBootstrapPeers() []ma.Multiaddr {
	peers := []ma.Multiaddr{}
	peers = append(peers, dht.DefaultBootstrapPeers...)
	sprawlBootstrapAddresses := []string{"/ip4/157.245.171.225/tcp/4001/ipfs/12D3KooWSNq7ujFYrMRJBKU51rJ9JvWr8tbzJ4e9cWTAC1TiXsfP"}
	for _, addr := range sprawlBootstrapAddresses {
		mAddr, _ := ma.NewMultiaddr(addr)
		peers = append(peers, mAddr)
	}
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
			if p2p.Logger != nil {
				p2p.Logger.Error(errors.E(errors.Op("Add dht"), err))
			}
		}
		return p2p.kademliaDHT, err
	}
	return libp2p.Routing(NewDHT)
}

// CreateOptions queries p2p.Config for any user-submitted options and assigns defaults
func (p2p *P2p) CreateOptions() []libp2pConfig.Option {
	options := []libp2pConfig.Option{}
	externalIP := p2p.Config.GetString("p2p.externalIP")
	p2pPort := p2p.Config.GetString("p2p.port")

	// Non-configurable options, since we always need an identity and the DHT discovery
	options = append(options, p2p.initDHT())
	options = append(options, libp2p.Identity(p2p.privateKey))

	// libp2p relay options
	if p2p.Config.GetBool("p2p.enableRelay") {
		options = append(options, libp2p.EnableRelay())
	}
	if p2p.Config.GetBool("p2p.enableAutoRelay") {
		options = append(options, libp2p.EnableAutoRelay())
	}

	// If NAT port map is not enabled, define listened addresses and port manually
	if p2p.Config.GetBool("p2p.enableNATPortMap") {
		options = append(options, libp2p.NATPortMap())
	} else {
		multiaddrs := defaultListenAddrs(p2pPort)
		if externalIP != "" {
			extMultiAddr, err := createMultiAddr(externalIP, p2pPort)
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

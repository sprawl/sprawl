package p2p

import (
	"fmt"

	libp2p "github.com/libp2p/go-libp2p"
	libp2pConfig "github.com/libp2p/go-libp2p/config"
	ma "github.com/multiformats/go-multiaddr"
)

func defaultListenAddrs(p2pPort string) []ma.Multiaddr {
	multiaddrs := []ma.Multiaddr{}
	local1, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", p2pPort))
	local2, _ := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%s", p2pPort))
	multiaddrs = append(multiaddrs, local1)
	multiaddrs = append(multiaddrs, local2)
	return multiaddrs
}

// CreateOptions queries p2p.Config for any user-submitted options and assigns defaults
func (p2p *P2p) CreateOptions() []libp2pConfig.Option {
	options := []libp2pConfig.Option{}
	externalIP := p2p.Config.GetString("p2p.externalIP")
	p2pPort := p2p.Config.GetString("p2p.port")

	options = append(options, p2p.initDHT())
	options = append(options, libp2p.Identity(p2p.privateKey))

	if p2p.Config.GetBool("p2p.enableRelay") {
		options = append(options, libp2p.EnableRelay())
	}

	if p2p.Config.GetBool("p2p.enableAutoRelay") {
		options = append(options, libp2p.EnableAutoRelay())
	}

	if p2p.Config.GetBool("p2p.enableNATPortMap") {
		options = append(options, libp2p.NATPortMap())
	} else {
		multiaddrs := defaultListenAddrs(p2pPort)
		if externalIP != "" {
			extMultiAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", externalIP, p2pPort))
			if err != nil {
				p2p.Logger.Errorf("Couldn't create multiaddr: %v", err)
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

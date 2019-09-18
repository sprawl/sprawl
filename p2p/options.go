package p2p

import (
	libp2p "github.com/libp2p/go-libp2p"
	libp2pConfig "github.com/libp2p/go-libp2p/config"
)

func (p2p *P2p) CreateOptions() []libp2pConfig.Option {
	options := []libp2pConfig.Option{}
	p2p.Logger.Infof("p1234 %v", p2p.kademliaDHT)
	if p2p.Config.GetBool("p2p.options.enableDHT") {
		options = append(options, p2p.initDHT())
	}
	p2p.Logger.Infof("p5678 %v", p2p.kademliaDHT)

	if p2p.Config.GetBool("p2p.options.enableIdentity") {
		options = append(options, libp2p.Identity(p2p.privateKey))
	}

	if p2p.Config.GetBool("p2p.options.enableRelay") {
		options = append(options, libp2p.EnableRelay())
	}

	if p2p.Config.GetBool("p2p.options.enableAutoRelay") {
		options = append(options, libp2p.EnableAutoRelay())
	}

	if p2p.Config.GetBool("p2p.options.enableNATPortMap") {
		options = append(options, libp2p.NATPortMap())
	}

	p2p.Logger.Infof("p9101112 %v", p2p.kademliaDHT)

	return options
}

package p2p

import (
	libp2p "github.com/libp2p/go-libp2p"
	config "github.com/libp2p/go-libp2p/config"
)

func (p2p *P2p)CreateOptions() []config.Option {
	options := []config.Option{}
	options = append(options, p2p.initDHT())
	options = append(options, libp2p.Identity(p2p.privateKey))
	options = append(options, libp2p.EnableRelay())
	options = append(options, libp2p.EnableAutoRelay())
	options = append(options, libp2p.NATPortMap())
	return options
}
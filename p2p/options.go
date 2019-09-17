package p2p

import (
	config "github.com/eqlabs/sprawl/config"
	libp2p "github.com/libp2p/go-libp2p"
	libp2pConfig "github.com/libp2p/go-libp2p/config"
)

var appConfig *config.Config

func init() {
	appConfig = &config.Config{}
	appConfig.ReadConfig("./config/default")
}

func (p2p *P2p) CreateOptions() []libp2pConfig.Option {
	options := []libp2pConfig.Option{}
	if appConfig.GetBool("p2p.options.enableDHT") {
		options = append(options, p2p.initDHT())
	}

	if appConfig.GetBool("p2p.options.enableIdentity") {
		options = append(options, libp2p.Identity(p2p.privateKey))
	}

	if appConfig.GetBool("p2p.options.enableRelay") {
		options = append(options, libp2p.EnableRelay())
	}

	if appConfig.GetBool("p2p.options.enableAutoRelay") {
		options = append(options, libp2p.EnableAutoRelay())
	}

	if appConfig.GetBool("p2p.options.enableNATPortMap") {
		options = append(options, libp2p.NATPortMap())
	}

	return options
}

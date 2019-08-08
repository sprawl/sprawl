package main

import (
	"github.com/eqlabs/sprawl/api"
	"github.com/eqlabs/sprawl/p2p"
)

func main() {
	p2pImpl := &p2p.P2p{}
	p2pImpl.Run()
	api.Run(1337, p2pImpl)
}

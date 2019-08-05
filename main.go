package main

import (
	"github.com/eqlabs/sprawl/api"
	"github.com/eqlabs/sprawl/p2p"
)

func main() {
	p2p.Run()
	api.Init()
}

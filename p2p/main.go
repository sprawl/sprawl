package p2p

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	network "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	multiaddr "github.com/multiformats/go-multiaddr"
)

type P2p struct {
	config           Config
	ctx              context.Context
	host             host.Host
	kademliaDHT      *dht.IpfsDHT
	routingDiscovery *discovery.RoutingDiscovery
	peerChan         <-chan peer.AddrInfo
}

func handleStream(stream network.Stream) {

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)

	// 'stream' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

func (p2p *P2p) createConfig() {
	var err error
	p2p.config, err = ParseFlags()
	if err != nil {
		panic(err)
	}
}

func (p2p *P2p) createContext() {
	p2p.ctx = context.Background()
}

func (p2p *P2p) createHost(ctx context.Context, config Config) {
	var err error
	p2p.host, err = libp2p.New(ctx,
		libp2p.ListenAddrs([]multiaddr.Multiaddr(config.ListenAddresses)...),
	)
	if err != nil {
		panic(err)
	}
}

func (p2p *P2p) createKademliaDHT(ctx context.Context, host host.Host) {
	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	var err error
	p2p.kademliaDHT, err = dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

}

func (p2p *P2p) bootstrapDHT(ctx context.Context, kademliaDHT *dht.IpfsDHT) {
	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	var err error
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}
}

func (p2p *P2p) getPeerAddresses(ctx context.Context, config Config, host host.Host) {
	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.
	var wg sync.WaitGroup
	for _, peerAddr := range config.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(peerinfo)
			}
		}()
	}
	wg.Wait()
}

func (p2p *P2p) createRoutingDiscovery(kademliaDHT *dht.IpfsDHT) {
	p2p.routingDiscovery = discovery.NewRoutingDiscovery(kademliaDHT)
}

func (p2p *P2p) advertise(ctx context.Context, config Config, routingDiscovery *discovery.RoutingDiscovery) {
	discovery.Advertise(ctx, routingDiscovery, config.RendezvousString)
}

func (p2p *P2p) findPeers(ctx context.Context, config Config, routingDiscovery *discovery.RoutingDiscovery) {
	var err error
	p2p.peerChan, err = routingDiscovery.FindPeers(ctx, config.RendezvousString)
	if err != nil {
		panic(err)
	}
}

func (p2p *P2p) communicateWithPeers(ctx context.Context, config Config, host host.Host, peerChan <-chan peer.AddrInfo) {
	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(config.ProtocolID))

		if err != nil {
			continue
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeData(rw)
			go readData(rw)
		}
	}
}

// Run runs the p2p network
func Run() {
	p2p := P2p{}
	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	p2p.createConfig()
	p2p.createContext()
	p2p.createHost(p2p.ctx, p2p.config)
	p2p.createKademliaDHT(p2p.ctx, p2p.host)
	p2p.host.SetStreamHandler(protocol.ID(p2p.config.ProtocolID), handleStream)
	p2p.bootstrapDHT(p2p.ctx, p2p.kademliaDHT)
	p2p.getPeerAddresses(p2p.ctx, p2p.config, p2p.host)
	p2p.createRoutingDiscovery(p2p.kademliaDHT)
	p2p.advertise(p2p.ctx, p2p.config, p2p.routingDiscovery)
	p2p.findPeers(p2p.ctx, p2p.config, p2p.routingDiscovery)
	p2p.communicateWithPeers(p2p.ctx, p2p.config, p2p.host, p2p.peerChan)
	select {}
}

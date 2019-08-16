package p2p

import (
	"bufio"
	"context"
	"fmt"
	"sync"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	network "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	configt "github.com/libp2p/go-libp2p/config"
	routing "github.com/libp2p/go-libp2p-core/routing"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	multiaddr "github.com/multiformats/go-multiaddr"
	// discovery "github.com/libp2p/go-libp2p/p2p/discovery"
)

const pubsubTopic = "/"
const pubsubTopics = "/ssss"

// func (p2p *P2p)addPeers(ps *pubsub.PubSub, config Config){
// 	for peer := range p2p.peerChan {
// 		if peer.ID == p2p.host.ID() {
// 			continue
// 		}
// 		ps.AddPeer(peer.ID, protocol.ID(config.ProtocolID))
// 		fmt.Println("SOmething happenedd!!!!!!!!!")
// 	}
// }

func publishaa(ps *pubsub.PubSub, ctx context.Context, hos host.Host){
	for {
		err := ps.Publish(pubsubTopic, []byte("MÄMMI"))
		fmt.Println("jouuuuuuussbost ")
		if err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
		fmt.Println("jouuuuuuussbostsadsadg %s", ps.ListPeers(pubsubTopic))
		fmt.Println("jouuuuuuussdg %s", hos.Addrs())
		fmt.Println("jouuuusstsadsadg %s", hos.Peerstore())
	}
}

func (p2p *P2p) pubsub(ctx context.Context, host host.Host) {
	fmt.Println("jouuuuuusdsdsdu")
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		panic(err)
	}
	sub, err := ps.Subscribe(pubsubTopic)
	if err != nil {
		panic(err)
	}
	go publishaa(ps, ctx, host)
	// go p2p.addPeers(ps, p2p.config)
	fmt.Println("Mammi")
	for {
		msg, _ := sub.Next(ctx)
		fmt.Println("TULEE VARMASTI %s", msg)
	}
}

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
	reader := bufio.NewReader(stream)

	go readData(reader, p2p.host)
	writer := bufio.NewWriter(stream)
	writeData(writer, []byte("ALa sano sitaasddsa!\n"))

	// 'stream' will stay open until you close it (or the other side closes it).
}

func readData(reader *bufio.Reader, hos host.Host) {
	for {
		fmt.Println("mitt'r")
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}
		if bytes == nil {
			fmt.Println("mitt'r")
			return
		}
		if bytes[0] != byte('\n') {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Println("\x1b[32m%s\x1b[0m> ", bytes)
			fmt.Println("\x1b[32m%s\x1b[0m> ", hos.Peerstore)
			fmt.Println("\x1b[32mb%s\x1b[0m> ", hos.Addrs())
		}
	}
}

func writeData(writer *bufio.Writer, input []byte) {
	fmt.Println("Testi132")
	_, err := writer.Write(input)
	if err != nil {
		fmt.Println("Error writing to buffer")
		panic(err)
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing buffer")
		panic(err)
	}
}

func (p2p *P2p) createConfig() {
	var err error
	p2p.config, err = ParseFlags()
	fmt.Println("%s", p2p.config.ListenAddresses)
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
				fmt.Println("Kekkonen + %s", peerinfo)
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

func (p2p *P2p) SendToPeers(input []byte) {
	p2p.sendToPeers(p2p.ctx, p2p.config, p2p.host, p2p.peerChan, input)
}

func (p2p *P2p) sendToPeers(ctx context.Context, config Config, host host.Host, peerChan <-chan peer.AddrInfo, input []byte) {
	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(config.ProtocolID))

		if err != nil {
			continue
		} else {
			writer := bufio.NewWriter(stream)
			writeData(writer, input)
		}
		}
	}


func (p2p *P2p) listenPeers(ctx context.Context, config Config, host host.Host, peerChan <-chan peer.AddrInfo) {

	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}
		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(config.ProtocolID))

		if err != nil {
			continue
		} else {
			reader := bufio.NewReader(stream)
			go readData(reader, host)
			writer := bufio.NewWriter(stream)
			writeData(writer, []byte("ALa sano sita!\n"))
		}
	}
}

//KAUNIS MIELI MOODI
func (p2p *P2p) tinamenkka() configt.Option {
	NewDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		p2p.kademliaDHT, err = dht.New(p2p.ctx, h)
		return p2p.kademliaDHT, err
	}
	return libp2p.Routing(NewDHT)

}

func (p2p *P2p) tinamämmi(routing configt.Option) {
	var err error
	p2p.host, err = libp2p.New(p2p.ctx,
		libp2p.ListenAddrs([]multiaddr.Multiaddr(p2p.config.ListenAddresses)...), routing,
	)
	if err != nil {
		panic(err)
	}
}

// Run runs the p2p network
func (p2p *P2p) Run() {
	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	p2p.createConfig()
	fmt.Println("1Loppu")
	p2p.createContext()
	fmt.Println("2Loppu")
	// p2p.createHost(p2p.ctx, p2p.config)
	fmt.Println("3Loppu")
	p2p.tinamämmi(p2p.tinamenkka())
	// p2p.createKademliaDHT(p2p.ctx, p2p.host)
	fmt.Println("4Loppu")
	p2p.host.SetStreamHandler(protocol.ID(p2p.config.ProtocolID), p2p.handleStream)
	fmt.Println("6Loppu")
	p2p.getPeerAddresses(p2p.ctx, p2p.config, p2p.host)
	fmt.Println("7Loppu")
	p2p.createRoutingDiscovery(p2p.kademliaDHT)
	fmt.Println("8Loppu")
	p2p.advertise(p2p.ctx, p2p.config, p2p.routingDiscovery)
	fmt.Println("9Loppu")
	p2p.findPeers(p2p.ctx, p2p.config, p2p.routingDiscovery)
	fmt.Println("10Loppu")
	fmt.Println("5Loppu")
	p2p.bootstrapDHT(p2p.ctx, p2p.kademliaDHT)
	p2p.pubsub(p2p.ctx, p2p.host)
	// fmt.Println("11Loppu")
	// p2p.listenPeers(p2p.ctx, p2p.config, p2p.host, p2p.peerChan)
	fmt.Println("12Loppu")
	select {}
}

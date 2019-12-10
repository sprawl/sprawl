package interfaces

import (
	"github.com/sprawl/sprawl/pb"
)

// P2p is a general p2p connection handler
type P2p interface {
	AddReceiver(receiver Receiver)
	Send(message *pb.WireMessage)
	Subscribe(channel *pb.Channel)
	Unsubscribe(channel *pb.Channel)
	GetAllPeers() []string
	BlacklistPeer(peerId *pb.Peer)
	Run()
	Close()
}

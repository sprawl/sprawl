package interfaces

import (
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/sprawl/sprawl/pb"
)

// P2p is a general p2p connection handler
type P2p interface {
	GetHostID() peer.ID
	GetHostIDString() string
	AddReceiver(receiver Receiver)
	Send(message *pb.WireMessage)
	Subscribe(channel *pb.Channel)
	Unsubscribe(channel *pb.Channel)
	GetAllPeers() []string
	BlacklistPeer(peerId *pb.Peer)
	OpenStream(peerID peer.ID) (Stream, error)
	CloseStream(peerID peer.ID) error
	Run()
	Close()
}

// Stream is a single stream instance between two peers
type Stream interface {
	WriteToStream(data []byte) error
}

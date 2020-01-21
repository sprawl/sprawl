package interfaces

import peer "github.com/libp2p/go-libp2p-core/peer"

// Receiver receives and parses all Wiremessages from p2p
type Receiver interface {
	Receive(data []byte, from peer.ID) error
}

package p2p

import (
	"context"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sprawl/sprawl/errors"
)

func (p2p *P2p) syncDataWithNewPeers(sub *pubsub.Subscription) {
	go func(ctx context.Context) {
		for {
			peerEvent, err := sub.NextPeerEvent(ctx)
			if p2p.Logger != nil {
				if !errors.IsEmpty(err) {
					p2p.Logger.Error(errors.E(errors.Op("Peer event"), err))
				}
			}
			if peerEvent.Type == 0 && peerEvent.Peer != p2p.host.ID() {
				if p2p.Logger != nil {
					p2p.Logger.Infof("New peer %s joined channel, syncing...", peerEvent.Peer)
				}
			}
		}
	}(p2p.ctx)
}

package p2p

import (
	"context"

	"github.com/golang/protobuf/proto"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/pb"
)

func (p2p *P2p) pingNewMembers(sub *pubsub.Subscription) {
	go func(ctx context.Context) {
		for {
			peerEvent, err := sub.NextPeerEvent(ctx)
			if p2p.Logger != nil {
				if !errors.IsEmpty(err) {
					p2p.Logger.Error(errors.E(errors.Op("Peer event"), err))
				}
			}

			p2p.Logger.Debugf("Peer event received from %s!", peerEvent.Peer)

			if peerEvent.Type == 0 && peerEvent.Peer.String() != p2p.host.ID().String() {
				p2p.Logger.Debug("Joined peer is not us!")

				from, err := peerEvent.Peer.Marshal()
				if !errors.IsEmpty(err) {
					if p2p.Logger != nil {
						p2p.Logger.Warn(errors.E(errors.Op("Marshal peerID in sync"), err))
					}
				}

				recipient := &pb.Recipient{PeerID: from}
				marshaledRecipient, err := proto.Marshal(recipient)
				if !errors.IsEmpty(err) {
					if p2p.Logger != nil {
						p2p.Logger.Warn(errors.E(errors.Op("Marshal recipient in sync"), err))
					}
				}

				marshaledSender, err := p2p.host.ID().Marshal()
				if !errors.IsEmpty(err) {
					if p2p.Logger != nil {
						p2p.Logger.Warn(errors.E(errors.Op("Marshal sender in sync"), err))
					}
				}

				wireMessage := &pb.WireMessage{ChannelID: []byte(sub.Topic()), Operation: pb.Operation_PING, Sender: marshaledSender, Data: marshaledRecipient}
				p2p.Send(wireMessage)
			}
		}
	}(p2p.ctx)
}

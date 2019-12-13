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
			if peerEvent.Type == 0 && peerEvent.Peer != p2p.host.ID() {
				recipient := &pb.Recipient{PeerID: []byte(peerEvent.Peer)}
				marshaledRecipient, err := proto.Marshal(recipient)
				if !errors.IsEmpty(err) {
					if p2p.Logger != nil {
						p2p.Logger.Warn(errors.E(errors.Op("Marshal recipient"), err))
					}
				}
				wireMessage := &pb.WireMessage{ChannelID: []byte(sub.Topic()), Operation: pb.Operation_PING, Data: marshaledRecipient}
				p2p.Send(wireMessage)
			}
		}
	}(p2p.ctx)
}

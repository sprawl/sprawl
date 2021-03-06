package p2p

import (
	"context"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/pb"
)

func (p2p *P2p) listenToChannel(ctx context.Context, sub *pubsub.Subscription, channel *pb.Channel) {
	go func(ctx context.Context) {
		for {
			msg, err := sub.Next(ctx)
			if !errors.IsEmpty(err) {
				p2p.Logger.Error(errors.E(errors.Op("Next Message"), err))
				return
			}

			data := msg.GetData()
			peer := msg.GetFrom()

			if peer != p2p.host.ID() {
				if p2p.Receiver != nil {
					err = p2p.Receiver.Receive(data, peer)
					if !errors.IsEmpty(err) {
						p2p.Logger.Error(errors.E(errors.Op("Receive data"), err))
					}
				} else {
					p2p.Logger.Warn("Receiver not registered with p2p, not parsing any incoming data!")
				}
			}
		}
	}(ctx)
}

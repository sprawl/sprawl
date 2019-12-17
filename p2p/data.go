package p2p

import (
	"context"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/pb"
)

func (p2p *P2p) listenToChannel(sub *pubsub.Subscription, channel *pb.Channel, quitSignal chan bool) {
	go func(ctx context.Context) {
		for {
			msg, err := sub.Next(ctx)
			if !errors.IsEmpty(err) {
				if p2p.Logger != nil {
					p2p.Logger.Error(errors.E(errors.Op("Next Message"), err))
				}
			}

			data := msg.GetData()
			peer := msg.GetFrom()

			if peer != p2p.host.ID() {
				if p2p.Receiver != nil {
					err = p2p.Receiver.Receive(data)
					if !errors.IsEmpty(err) {
						if p2p.Logger != nil {
							p2p.Logger.Error(errors.E(errors.Op("Receive data"), err))
						}
					}
				} else {
					if p2p.Logger != nil {
						p2p.Logger.Warn("Receiver not registered with p2p, not parsing any incoming data!")
					}
				}
			}

			select {
			case quit := <-quitSignal: //Delete subscription
				if quit {
					delete(p2p.subscriptions, string(channel.GetId()))
					return
				}
			default:
			}
		}
	}(p2p.ctx)
}

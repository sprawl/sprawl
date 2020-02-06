package p2p

import (
	"context"

	"github.com/golang/protobuf/proto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/pb"
)

func (p2p *P2p) requestSync(ctx context.Context, topicString string, topic *pubsub.Topic) {
	eventHandler, err := topic.EventHandler()
	if !errors.IsEmpty(err) {
		p2p.Logger.Error(errors.E(errors.Op("Return topic's event handler"), err))
	}

	//Add alternative if this fail
	go func(ctx context.Context) {
		for {
		peerEvent, err := eventHandler.NextPeerEvent(ctx)
		if !errors.IsEmpty(err) {
			p2p.Logger.Error(errors.E(errors.Op("Get next peer event"), err))
		}
		if peerEvent.Type == 0 && peerEvent.Peer.String() != p2p.host.ID().String() {
			p2p.sendSyncRequest(peerEvent.Peer, topicString)
			if !errors.IsEmpty(err) {
				p2p.Logger.Error(errors.E(errors.Op("Request sync"), err))
				} else {
					break
				}
			}
		}
	}(ctx)
}

func (p2p *P2p) sendSyncRequest(peerID peer.ID, topicString string) error {
	stream, err := p2p.OpenStream(peerID)
	if !errors.IsEmpty(err) {
		return errors.E(errors.Op("Open a sync stream"), err)
	}
	syncMessage := &pb.WireMessage{Operation: pb.Operation_SYNC_REQUEST, ChannelID: []byte(topicString), Data: nil}

	marshaledData, err := proto.Marshal(syncMessage)
	if !errors.IsEmpty(err) {
		return errors.E(errors.Op("Marshal sync request wireMessage"), err)
	}
	err = stream.WriteToStream(marshaledData)
	if !errors.IsEmpty(err) {
		return errors.E(errors.Op("Write sync request to stream"), err)
	}
	err = p2p.CloseStream(peerID)
	if !errors.IsEmpty(err) {
		return errors.E(errors.Op("Close the stream"), err)
	}
	return nil
}

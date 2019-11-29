package interfaces

import (
	"github.com/sprawl/sprawl/pb"
)

type P2p interface {
	RegisterOrderService(orders OrderService)
	RegisterChannelService(channels ChannelService)
	Send(message *pb.WireMessage)
	Subscribe(channel *pb.Channel)
	Unsubscribe(channel *pb.Channel)
	Run()
}

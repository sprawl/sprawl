package interfaces

import (
	"github.com/eqlabs/sprawl/pb"
)

type P2p interface {
	RegisterOrderService(orders OrderService)
	RegisterChannelService(channels ChannelService)
	Input(data []byte, channel *pb.Channel)
	Subscribe(channel *pb.Channel)
	Run()
}

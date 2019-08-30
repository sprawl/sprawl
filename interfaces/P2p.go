package interfaces

import (
	"github.com/eqlabs/sprawl/pb"
)

type P2p interface {
	RegisterOrderService(orders OrderService)
	RegisterChannelService(channels ChannelService)
	Input(channel pb.Channel, data []byte)
	Subscribe(channel pb.Channel)
	Run()
}

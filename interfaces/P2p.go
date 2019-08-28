package interfaces

type P2p interface {
	RegisterOrderService(orders OrderService)
	RegisterChannelService(channels ChannelService)
	Input(order []byte, channel string)
	Subscribe(channel string)
	Run()
}

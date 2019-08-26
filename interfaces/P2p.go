package interfaces

type P2p interface {
	RegisterOrderService(orders OrderService)
	RegisterChannelService(channels ChannelService)
	Input(data []byte, topic string)
	Subscribe(topic string)
	Run()
}

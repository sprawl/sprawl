package interfaces

type P2p interface {
	RegisterOrderService(orders OrderService)
	RegisterChannelService(channels ChannelService)
	PublishMessage(topic string, input []byte)
	Subscribe(topic string)
	Run()
}

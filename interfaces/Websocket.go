package interfaces

import "github.com/sprawl/sprawl/pb"

type WebsocketService interface {
	Start()
	Close()
	PushToWebsockets(message *pb.WireMessage)
}

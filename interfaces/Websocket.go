package interfaces

import "github.com/sprawl/sprawl/pb"

type WebsocketService interface {
	Start()
	Close()
	RelayToClients(message *pb.WireMessage)
}

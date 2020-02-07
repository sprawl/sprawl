package interfaces

import (
	"context"

	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/sprawl/sprawl/pb"
)

// OrderService is an interface to the Order endpoints in sprawl.proto
type OrderService interface {
	RegisterStorage(db Storage)
	RegisterP2p(p2p P2p)
	RegisterWebsocket(websocket WebsocketService)
	Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error)
	Receive(data []byte, from peer.ID) error
	Delete(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.Empty, error)
	Lock(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.Empty, error)
	Unlock(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.Empty, error)
	GetOrder(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.Order, error)
	GetAllOrders(ctx context.Context, in *pb.Empty) (*pb.OrderList, error)
}

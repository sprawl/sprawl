package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"

	"github.com/eqlabs/sprawl/db"
	"github.com/golang/protobuf/proto"
	ptypes "github.com/golang/protobuf/ptypes"
)

var storage = &db.Storage{}

func init() {
	// Initialize storage
	storage.SetDbPath("/var/lib/sprawl/data")
	storage.Run()
}

// OrderService implements the OrderService Server service.proto
type OrderService struct {
	channels []Channel
}

// Create creates an Order, storing it locally and broadcasts the Order to all other nodes on the channel
func (s *OrderService) Create(ctx context.Context, in *CreateRequest) (*CreateResponse, error) {
	// Get current timestamp as protobuf type
	now := ptypes.TimestampNow()

	// TODO: Use the node's private key here as a secret to sign the Order ID with
	secret := "mysecret"

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write(append([]byte(in.String()), []byte(now.String())...))

	// Get result and encode as hexadecimal string
	id := h.Sum(nil)

	// Construct the order
	order := &Order{
		Id:           id,
		Created:      now,
		Asset:        in.Asset,
		CounterAsset: in.CounterAsset,
		Amount:       in.Amount,
		Price:        in.Price,
		State:        State_OPEN,
	}

	// Get order as bytes
	orderInBytes, err := proto.Marshal(order)
	if err != nil {
		panic(err)
	}

	// Save order to LevelDB locally
	err = storage.Put(id, orderInBytes)
	if err != nil {
		panic(err)
	}

	// TODO: Propagate order to other nodes via sprawl/p2p

	// TODO: Properly return any errors to client instead of panicking
	// Return the response to the gRPC client
	return &CreateResponse{
		CreatedOrder: order,
		Error:        nil,
	}, nil
}

// Delete removes the Order with the specified ID locally, and broadcasts the same request to all other nodes on the channel
func (s *OrderService) Delete(ctx context.Context, in *OrderSpecificRequest) (*GenericResponse, error) {
	// Try to delete the Order from LevelDB with specified ID
	err := storage.Delete(in.GetId())
	if err != nil {
		panic(err)
	}

	// TODO: Propagate the deletion to other nodes via sprawl/p2p
	// TODO: Properly return any errors to client instead of panicking
	return &GenericResponse{
		Error: nil,
	}, nil
}

// Lock locks the given Order if the Order is created by this node, broadcasts the lock to other nodes on the channel.
func (s *OrderService) Lock(ctx context.Context, in *OrderSpecificRequest) (*GenericResponse, error) {

	// TODO: Add Order locking logic

	return &GenericResponse{
		Error: nil,
	}, nil
}

// Unlock unlocks the given Order if it's created by this node, broadcasts the unlocking operation to other nodes on the channel.
func (s *OrderService) Unlock(ctx context.Context, in *OrderSpecificRequest) (*GenericResponse, error) {

	// TODO: Add Order unlocking logic

	return &GenericResponse{
		Error: nil,
	}, nil
}

// Join joins a channel, starting a new instance of libp2p in OrderService.channels
func (s *OrderService) Join(ctx context.Context, in *Channel) (*JoinResponse, error) {

	// TODO: Add Channel joining logic

	return &JoinResponse{
		JoinedChannel: &Channel{},
	}, nil
}

// Leave leaves a channel, removing an instance of libp2p from OrderService.channelsi
func (s *OrderService) Leave(ctx context.Context, in *Channel) (*GenericResponse, error) {

	// TODO: Add Channel leaving logic

	return &GenericResponse{
		Error: nil,
	}, nil
}

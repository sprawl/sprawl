package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"

	ptypes "github.com/golang/protobuf/ptypes"
)

// OrderService implements the OrderService Server service.proto
type OrderService struct {
	channels []Channel
}

// Create creates an Order, storing it locally and broadcasts the Order to all other nodes on the channel
func (s *OrderService) Create(ctx context.Context, in *CreateRequest) (*CreateResponse, error) {

	// TODO: Add Order creation logic, save & propagate
	now := ptypes.TimestampNow()

	secret := "mysecret"

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write(append([]byte(in.String()), []byte(now.String())...))

	// Get result and encode as hexadecimal string
	id := h.Sum(nil)

	order := &Order{
		Id:           id,
		Created:      now,
		Asset:        in.Asset,
		CounterAsset: in.CounterAsset,
		Amount:       in.Amount,
		Price:        in.Price,
		State:        State_OPEN,
	}

	return &CreateResponse{
		CreatedOrder: order,
		Error:        nil,
	}, nil
}

// Delete removes the Order with the specified ID locally, and broadcasts the same request to all other nodes on the channel
func (s *OrderService) Delete(ctx context.Context, in *OrderSpecificRequest) (*GenericResponse, error) {

	// TODO: Add order deletion logic

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

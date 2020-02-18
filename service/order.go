package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"strings"

	"github.com/golang/protobuf/proto"
	ptypes "github.com/golang/protobuf/ptypes"
	"github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/identity"
	"github.com/sprawl/sprawl/interfaces"
	"github.com/sprawl/sprawl/pb"
)

// OrderService implements the OrderService Server service.proto
type OrderService struct {
	Logger    interfaces.Logger
	Storage   interfaces.Storage
	P2p       interfaces.P2p
	websocket interfaces.WebsocketService
}

func getOrderStorageKey(channelID []byte, orderID []byte) []byte {
	return []byte(strings.Join([]string{string(interfaces.OrderPrefix), string(channelID), string(orderID)}, ""))
}

func getOrderQueryPrefix(channelID []byte) []byte {
	return []byte(strings.Join([]string{string(interfaces.OrderPrefix), string(channelID)}, ""))
}

// RegisterWebsocket registers a websocket service to enable websocket connections between client and node
func (s *OrderService) RegisterWebsocket(websocket interfaces.WebsocketService) {
	s.websocket = websocket
}

// RegisterStorage registers a storage service to store the Orders in
func (s *OrderService) RegisterStorage(storage interfaces.Storage) {
	s.Storage = storage
}

// RegisterP2p registers a p2p service
func (s *OrderService) RegisterP2p(p2p interfaces.P2p) {
	s.P2p = p2p
}

// GetSignature generates signature from order and returns it
func (s *OrderService) GetSignature(order *pb.Order) ([]byte, error) {
	orderCopy := *order
	orderCopy.State = pb.State_OPEN
	orderCopy.Signature = nil
	orderCopy.Nonce = 0
	orderInBytes, err := proto.Marshal(&orderCopy)
	if !errors.IsEmpty(err) {
		return nil, errors.E(errors.Op("Marshal order in GetSignature"), err)
	}

	return identity.Sign(s.Storage, orderInBytes)
}

// VerifyOrder verifies order
func (s *OrderService) VerifyOrder(publicKey crypto.PubKey, order *pb.Order) (bool, error) {
	orderCopy := *order
	sig := orderCopy.Signature
	orderCopy.Signature = nil
	orderCopy.State = pb.State_OPEN
	orderCopy.Nonce = 0
	orderInBytes, err := proto.Marshal(&orderCopy)
	if !errors.IsEmpty(err) {
		return false, errors.E(errors.Op("Marshal order in VerifyOrder"), err)
	}
	return identity.Verify(publicKey, orderInBytes, sig)
}

// Create creates an Order, storing it locally and broadcasts the Order to all other nodes on the channel
func (s *OrderService) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {

	_, publicKey, err := identity.GetIdentity(s.Storage)
	if !errors.IsEmpty(err) {
		errors.E(errors.Op("Get public key in create order"), err)
	}

	// Get current timestamp as protobuf type
	now := ptypes.TimestampNow()

	secret, err := publicKey.Bytes()
	if !errors.IsEmpty(err) {
		errors.E(errors.Op("Turn public key into bytes"), err)
	}

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, secret)

	// Write Data to it
	h.Write(append([]byte(in.String()), []byte(now.String())...))

	// Get result and encode as hexadecimal string
	id := h.Sum(nil)

	// Construct the order
	order := &pb.Order{
		Id:           id,
		Created:      now,
		Asset:        in.Asset,
		CounterAsset: in.CounterAsset,
		Amount:       in.Amount,
		Price:        in.Price,
		State:        pb.State_OPEN, //Mutable
		Nonce:        0,             //Mutable
	}

	sig, err := s.GetSignature(order)
	if !errors.IsEmpty(err) {
		return &pb.CreateResponse{
			CreatedOrder: order,
		}, errors.E(errors.Op("Get Signature"), err)
	}

	order.Signature = sig

	// Get order as bytes
	orderInBytes, err := proto.Marshal(order)
	if !errors.IsEmpty(err) {
		s.Logger.Warn(errors.E(errors.Op("Marshal order"), err))
	}

	// Save order to LevelDB locally
	err = s.Storage.Put(getOrderStorageKey(in.GetChannelID(), id), orderInBytes)
	if !errors.IsEmpty(err) {
		err = errors.E(errors.Op("Put order"), err)
	}

	// Construct the message to send to other peers
	wireMessage := &pb.WireMessage{ChannelID: in.GetChannelID(), Operation: pb.Operation_CREATE, Data: orderInBytes}

	if s.P2p != nil {
		// Send the order creation by wire
		s.P2p.Send(wireMessage)
	} else {
		s.Logger.Warn("P2p service not registered with OrderService, not publishing or receiving orders from the network!")
	}

	return &pb.CreateResponse{
		CreatedOrder: order,
	}, err
}

// Receive receives a buffer from p2p and tries to unmarshal it into a struct
func (s *OrderService) Receive(buf []byte, from peer.ID) error {
	wireMessage := &pb.WireMessage{}
	err := proto.Unmarshal(buf, wireMessage)
	if !errors.IsEmpty(err) {
		return errors.E(errors.Op("Unmarshal wiremessage proto in Receive"), err)
	}
	if s.websocket != nil {
		s.websocket.PushToWebsockets(wireMessage)
	}

	// Read operation and data from the WireMessage
	op := wireMessage.GetOperation()
	data := wireMessage.GetData()
	channelID := wireMessage.GetChannelID()

	s.Logger.Debugf("%s: %s.%s", from.String(), channelID, op)

	if s.Storage != nil {
		switch op {

		case pb.Operation_CREATE:
			// Validate order
			order := &pb.Order{}
			err = proto.Unmarshal(data, order)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Unmarshal order proto in Receive"), err)
			}

			publickey, err := from.ExtractPublicKey()
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Extract public key in Receive"), err)
			}
			isCreator, err := s.VerifyOrder(publickey, order)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Verify order creator in Receive"), err)
			}
			if isCreator {
				// Save order to LevelDB locally
				err = s.Storage.Put(getOrderStorageKey(channelID, order.GetId()), data)
				if !errors.IsEmpty(err) {
					err = errors.E(errors.Op("Put order"), err)
				}
			} else {
				s.Logger.Debug("Received create request from someone that doesn't own the order")
			}

		case pb.Operation_DELETE:
			// Unmarshal order to get its key, validate
			order := &pb.Order{}
			err = proto.Unmarshal(data, order)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Unmarshal order proto in Receive"), err)
			}
			publickey, err := from.ExtractPublicKey()
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Extract public key in Receive"), err)
			}

			isCreator, err := s.VerifyOrder(publickey, order)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Verify order creator in Receive"), err)
			}
			if isCreator {
				err = s.Storage.Delete(getOrderStorageKey(channelID, order.GetId()))
				if !errors.IsEmpty(err) {
					return errors.E(errors.Op("Delete order"), err)
				}
			} else {
				s.Logger.Debug("Received delete request from someone that doesn't own the order")
			}

		case pb.Operation_SYNC_REQUEST:
			orders, err := s.Storage.GetAllWithPrefix(string(getOrderQueryPrefix(channelID)))
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Fetch orders for sync"), err)
			}

			orderList := &pb.OrderList{}
			for _, value := range orders {
				order := &pb.Order{}
				proto.Unmarshal([]byte(value), order)
				orderList.Orders = append(orderList.Orders, order)
			}

			marshaledOrderList, err := proto.Marshal(orderList)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Marshal orderList in sync request"), err)
			}

			syncMessage := &pb.WireMessage{Operation: pb.Operation_SYNC_RECEIVE, ChannelID: channelID, Data: marshaledOrderList}

			marshaledData, err := proto.Marshal(syncMessage)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Marshal wireMessage in sync request"), err)
			}

			stream, err := s.P2p.OpenStream(from)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Open a sync request stream"), err)
			}

			err = stream.WriteToStream(marshaledData)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Write to stream"), err)
			}
			err = s.P2p.CloseStream(from)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Close the stream"), err)
			}

		case pb.Operation_SYNC_RECEIVE:
			orderList := &pb.OrderList{}
			err = proto.Unmarshal(data, orderList)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Unmarshal order proto in Receive"), err)
			}
			s.Logger.Info(orderList)
			for _, order := range orderList.GetOrders() {
				orderBytes, err := proto.Marshal(order)
				if !errors.IsEmpty(err) {
					err = errors.E(errors.Op("Marshal order from received orderList"), err)
				}
				err = s.Storage.Put(getOrderStorageKey(channelID, order.GetId()), orderBytes)
				if !errors.IsEmpty(err) {
					err = errors.E(errors.Op("Put order"), err)
				}
			}
		case pb.Operation_LOCK, pb.Operation_UNLOCK:
			// Unmarshal order to get its key, validate
			order := &pb.Order{}
			err = proto.Unmarshal(data, order)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Unmarshal order proto in Receive"), err)
			}

			previousOrderData, err := s.Storage.Get(getOrderStorageKey(channelID, order.GetId()))
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Get previous order"), err)
			}
			previousOrder := &pb.Order{}
			proto.Unmarshal(previousOrderData, previousOrder)
			if previousOrder.Nonce > order.Nonce {
				return errors.E(errors.Op("Compare nonces"), "new order is older")
			}

			publickey, err := from.ExtractPublicKey()
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Extract public key in Receive"), err)
			}

			isCreator, err := s.VerifyOrder(publickey, order)
			if !errors.IsEmpty(err) {
				return errors.E(errors.Op("Verify order creator in Receive"), err)
			}

			if isCreator {
				// Save order to LevelDB locally
				err = s.Storage.Put(getOrderStorageKey(channelID, order.GetId()), data)
				if !errors.IsEmpty(err) {
					return errors.E(errors.Op("Store lock/unlock order"), err)
				}
			} else {
				s.Logger.Debug("Received delete request from someone that doesn't own the order")
			}

		}
	} else {
		s.Logger.Warn("Storage not registered with OrderService, not persisting Orders!")
	}

	return err
}

// GetOrder fetches a single order from the database
func (s *OrderService) GetOrder(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.Order, error) {
	data, err := s.Storage.Get(getOrderStorageKey(in.GetChannelID(), in.GetOrderID()))
	if !errors.IsEmpty(err) {
		return nil, errors.E(errors.Op("Get order"), err)
	}
	order := &pb.Order{}
	proto.Unmarshal(data, order)
	return order, nil
}

// GetAllOrders fetches all orders from the database
func (s *OrderService) GetAllOrders(ctx context.Context, in *pb.Empty) (*pb.OrderList, error) {
	data, err := s.Storage.GetAllWithPrefix(string(interfaces.OrderPrefix))
	if !errors.IsEmpty(err) {
		return nil, errors.E(errors.Op("Get all orders"), err)
	}

	orders := make([]*pb.Order, 0)
	i := 0
	for _, value := range data {
		order := &pb.Order{}
		proto.Unmarshal([]byte(value), order)
		orders = append(orders, order)
		i++
	}

	OrderList := &pb.OrderList{Orders: orders}
	return OrderList, nil
}

// Delete removes the Order with the specified ID locally, and broadcasts the same request to all other nodes on the channel
func (s *OrderService) Delete(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.Empty, error) {
	orderInBytes, err := s.Storage.Get(getOrderStorageKey(in.GetChannelID(), in.GetOrderID()))
	if !errors.IsEmpty(err) {
		return nil, errors.E(errors.Op("Delete order"), err)
	}

	order := &pb.Order{}
	err = proto.Unmarshal(orderInBytes, order)
	if !errors.IsEmpty(err) {
		return &pb.Empty{}, errors.E(errors.Op("Unmarshal order proto in Delete"), err)
	}

	_, publickey, err := identity.GetIdentity(s.Storage)
	if !errors.IsEmpty(err) {
		return &pb.Empty{}, errors.E(errors.Op("Get public key in Delete"), err)
	}

	isCreator, err := s.VerifyOrder(publickey, order)
	if !errors.IsEmpty(err) {
		return &pb.Empty{}, errors.E(errors.Op("Verify the order"), err)
	}

	// Construct the message to send to other peers
	wireMessage := &pb.WireMessage{ChannelID: in.GetChannelID(), Operation: pb.Operation_DELETE, Data: orderInBytes}

	if s.P2p != nil {
		if isCreator {
			// Send the order creation by wire
			s.P2p.Send(wireMessage)
		}
	} else {
		s.Logger.Warn("P2p service not registered with OrderService, not publishing or receiving orders from the network!")
	}

	// Try to delete the Order from LevelDB with specified ID
	err = s.Storage.Delete(getOrderStorageKey(in.GetChannelID(), in.GetOrderID()))
	if !errors.IsEmpty(err) {
		err = errors.E(errors.Op("Delete order"), err)
	}

	return &pb.Empty{}, err
}

// Lock locks the given Order if the Order is created by this node, broadcasts the lock to other nodes on the channel.
func (s *OrderService) Lock(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.Empty, error) {

	orderInBytes, err := s.Storage.Get(getOrderStorageKey(in.GetChannelID(), in.GetOrderID()))
	if !errors.IsEmpty(err) {
		return nil, errors.E(errors.Op("Get order in Lock"), err)
	}

	order := &pb.Order{}
	err = proto.Unmarshal(orderInBytes, order)
	if !errors.IsEmpty(err) {
		return &pb.Empty{}, errors.E(errors.Op("Unmarshal order proto in Lock"), err)
	}

	if order.State == pb.State_LOCKED {
		return &pb.Empty{}, errors.E(errors.Op("Check state"), "Trying to lock something that is already locked")
	}

	_, publickey, err := identity.GetIdentity(s.Storage)
	if !errors.IsEmpty(err) {
		return &pb.Empty{}, errors.E(errors.Op("Get public key in Lock"), err)
	}

	isCreator, err := s.VerifyOrder(publickey, order)
	if !errors.IsEmpty(err) {
		return &pb.Empty{}, errors.E(errors.Op("Verify the order in Lock"), err)
	}

	order.State = pb.State_LOCKED

	// Get order as bytes
	orderInBytes, err = proto.Marshal(order)
	if !errors.IsEmpty(err) {
		s.Logger.Warn(errors.E(errors.Op("Marshal order"), err))
	}

	// Construct the message to send to other peers
	wireMessage := &pb.WireMessage{ChannelID: in.GetChannelID(), Operation: pb.Operation_LOCK, Data: orderInBytes}

	if s.P2p != nil {
		if isCreator {
			// Send the order creation by wire
			s.P2p.Send(wireMessage)
		}
	} else {
		s.Logger.Warn("P2p service not registered with OrderService, not publishing or receiving orders from the network!")
	}

	// Save order to LevelDB locally
	err = s.Storage.Put(getOrderStorageKey(in.GetChannelID(), in.GetOrderID()), orderInBytes)
	if !errors.IsEmpty(err) {
		err = errors.E(errors.Op("Put order"), err)
	}

	return &pb.Empty{}, nil
}

// Unlock unlocks the given Order if it's created by this node, broadcasts the unlocking operation to other nodes on the channel.
func (s *OrderService) Unlock(ctx context.Context, in *pb.OrderSpecificRequest) (*pb.Empty, error) {

	// TODO: Add Order unlocking logic

	return &pb.Empty{}, nil
}

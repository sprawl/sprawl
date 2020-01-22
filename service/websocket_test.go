package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/pb"
	"github.com/stretchr/testify/assert"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var testChannel *pb.Channel = &pb.Channel{Id: []byte("testChannel")}
var testOrder *pb.Order = &pb.Order{Asset: string("ETH"), CounterAsset: string("BTC"), Amount: 52152, Price: 0.2, Id: []byte("jgkahgkjal")}
var testOrderInBytes []byte
var testWireMessage *pb.WireMessage

func StartMockServer(websocketService *WebsocketService) (ws *websocket.Conn, s *httptest.Server, err error) {
	s = httptest.NewServer(http.HandlerFunc(websocketService.connect))
	defer func() {
		if err != nil {
			s.Close()
			if ws != nil {
				ws.Close()
			}
		}
	}()
	//Get websocket address from http
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	//Connect to the server
	ws, _, err = websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		err = errors.E(errors.Op("Dial to websocket"), err)
	}
	return
}

func CloseMockServer(ws *websocket.Conn, s *httptest.Server) error {
	s.Close()
	return ws.Close()
}

func TestConnectAndRelay(t *testing.T) {
	wss := WebsocketService{Logger: log}
	ws, s, err := StartMockServer(&wss)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, CloseMockServer(ws, s))
	}()
	assert.NoError(t, err)
	testOrderInBytes, err := proto.Marshal(testOrder)
	assert.NoError(t, err)
	testWireMessage = &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}
	wss.RelayToClients(testWireMessage)
	_, p, err := ws.ReadMessage()
	assert.NoError(t, err)
	testWireMessage2 := &pb.WireMessage{}
	proto.Unmarshal(p, testWireMessage2)
	testOrder2 := &pb.Order{}
	proto.Unmarshal(testWireMessage2.GetData(), testOrder2)
	assert.Equal(t, testWireMessage.GetData(), testWireMessage2.GetData())
	assert.Equal(t, testOrder.GetId(), testOrder2.GetId())

}

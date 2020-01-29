package service

import (
	"fmt"
	"net/url"
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

const port uint = 3000

func StartServer(websocketService *WebsocketService) (ws *websocket.Conn, err error) {
	go websocketService.Start()
	u := url.URL{Scheme: "ws", Host: "localhost:" + fmt.Sprint(port), Path: "/"}
	ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		err = errors.E(errors.Op("Dial to websocket"), err)
	}
	return
}

func TestConnectAndRelay(t *testing.T) {
	wss := WebsocketService{Logger: log, Port: port}
	ws, err := StartServer(&wss)
	defer wss.Close()
	assert.NoError(t, err)
	testOrderInBytes, err := proto.Marshal(testOrder)
	assert.NoError(t, err)
	testWireMessage = &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}
	wss.PushToWebsockets(testWireMessage)
	_, p, err := ws.ReadMessage()
	assert.NoError(t, err)
	testWireMessage2 := &pb.WireMessage{}
	proto.Unmarshal(p, testWireMessage2)
	testOrder2 := &pb.Order{}
	proto.Unmarshal(testWireMessage2.GetData(), testOrder2)
	assert.Equal(t, testWireMessage.GetData(), testWireMessage2.GetData())
	assert.Equal(t, testOrder.GetId(), testOrder2.GetId())

}

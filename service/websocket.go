package service

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/interfaces"
	"github.com/sprawl/sprawl/pb"
)

type WebsocketService struct {
	Connections []*websocket.Conn
	Logger      interfaces.Logger
	Port        uint
	httpServer  http.Server
}


func (ws *WebsocketService) connect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if !errors.IsEmpty(err) {
		if ws.Logger != nil {
			ws.Logger.Warn(errors.E(errors.Op("Upgrade from http to ws"), err))
		}
		return
	}
	ws.Connections = append(ws.Connections, conn)
}

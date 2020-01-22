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

func (ws *WebsocketService) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws.connect(w, r)
	})
	ws.httpServer = http.Server{Addr: ":" + string(ws.Port), Handler: mux}
	err := ws.httpServer.ListenAndServe()
	if !errors.IsEmpty(err) {
		if ws.Logger != nil {
			ws.Logger.Error(errors.E(errors.Op("Listen and serve port :"+string(ws.Port)), err))
		}
	}
}

func (ws *WebsocketService) Close() {
	err := ws.httpServer.Close()
	if !errors.IsEmpty(err) {
		if ws.Logger != nil {
			ws.Logger.Error(errors.E(errors.Op("Close http server")), err)
		}
	}
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

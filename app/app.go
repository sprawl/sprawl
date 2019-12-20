package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sprawl/sprawl/database/inmemory"
	"github.com/sprawl/sprawl/database/leveldb"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/identity"
	"github.com/sprawl/sprawl/interfaces"
	"github.com/sprawl/sprawl/p2p"
	"github.com/sprawl/sprawl/pb"
	"github.com/sprawl/sprawl/service"
)

// App ties Sprawl's services together
type App struct {
	Storage interfaces.Storage
	P2p     *p2p.P2p
	Server  *service.Server
	Logger  interfaces.Logger
	config  interfaces.Config
}

func (app *App) debugPinger() {
	var testChannel *pb.Channel = &pb.Channel{Id: []byte("testChannel")}
	app.P2p.Subscribe(testChannel)
	testRequest := &pb.CreateRequest{ChannelID: testChannel.GetId(), Asset: string("ETH"), CounterAsset: string("BTC"), Amount: 52153, Price: 0.2}

	for {
		if app.Logger != nil {
			app.Logger.Infof("Debug pinger is sending testRequest: %s\n", testRequest)
		}
		orderID, err := app.Server.Orders.Create(context.Background(), testRequest)
		if !errors.IsEmpty(err) && app.Logger != nil {
			app.Logger.Error(errors.E(errors.Op("Create Request"), err))
		}
		testOrderSpecificRequest := &pb.OrderSpecificRequest{OrderID: orderID.GetCreatedOrder().GetId(), ChannelID: testChannel.GetId()}
		time.Sleep(time.Minute)
		app.Server.Orders.Delete(context.Background(), testOrderSpecificRequest)
	}
}

// InitServices ties the services together before running
func (app *App) InitServices(config interfaces.Config, Logger interfaces.Logger) {
	app.config = config
	app.Logger = Logger
	errors.SetDebug(app.config.GetBool("errors.enableStackTrace"))

	if app.Logger != nil {
		app.Logger.Infof("Saving data to %s", app.config.GetString("database.path"))
	}

	systemSignals := make(chan os.Signal)
	signal.Notify(systemSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-systemSignals:
			app.Logger.Infof("Received %s signal, shutting down.\n", sig)
			app.Server.Close()
			app.Storage.Close()
			app.P2p.Close()
			os.Exit(0)
		}
	}()

	// Start up the database
	if app.config.GetBool("database.inMemory") {
		app.Storage = &inmemory.Storage{
			Db: make(map[string]string),
		}
	} else {
		app.Storage = &leveldb.Storage{}
	}
	app.Storage.SetDbPath(app.config.GetString("database.path"))
	app.Storage.Run()

	privateKey, publicKey, err := identity.GetIdentity(app.Storage)

	if !errors.IsEmpty(err) && app.Logger != nil {
		app.Logger.Error(errors.E(errors.Op("Get identity"), err))
	}

	// Run the P2P process
	app.P2p = p2p.NewP2p(Logger, config, privateKey, publicKey)

	// Construct the server struct
	app.Server = service.NewServer(Logger, app.Storage, app.P2p)

	// Connect the order and channel services with p2p
	app.P2p.RegisterOrderService(app.Server.Orders)
	app.P2p.RegisterChannelService(app.Server.Channels)

	// Run the P2p service before running the gRPC server
	app.P2p.Run()
}

// Run is a separated main-function to ease testing
func (app *App) Run() {
	// Run the gRPC API

	defer app.Server.Close()
	defer app.Storage.Close()
	defer app.P2p.Close()

	if app.config.GetBool("p2p.debug") {
		if app.Logger != nil {
			app.Logger.Info("Running the debug pinger on channel \"testChannel\"!")
		}
		go app.debugPinger()
	}
	app.Server.Run(app.config.GetUint("rpc.port"))
}

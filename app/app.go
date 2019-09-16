package app

import (
	"time"

	"github.com/eqlabs/sprawl/db"
	"github.com/eqlabs/sprawl/identity"
	"github.com/eqlabs/sprawl/interfaces"
	"github.com/eqlabs/sprawl/p2p"
	"github.com/eqlabs/sprawl/pb"
	"github.com/eqlabs/sprawl/service"
	"github.com/gogo/protobuf/proto"
)

// App ties Sprawl's services together
type App struct {
	Storage *db.Storage
	P2p     *p2p.P2p
	Server  *service.Server
	log     interfaces.Logger
	config  interfaces.Config
}

func (app *App) debugPinger() {
	var testChannel *pb.Channel = &pb.Channel{Id: []byte("testChannel")}
	app.P2p.Subscribe(testChannel)

	var testOrder *pb.Order = &pb.Order{Asset: string("ETH"), CounterAsset: string("BTC"), Amount: 52152, Price: 0.2, Id: []byte("jgkahgkjal")}
	testOrderInBytes, err := proto.Marshal(testOrder)
	if err != nil {
		panic(err)
	}

	testWireMessage := &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}

	for {
		app.log.Infof("Debug pinger is sending testWireMessage: %s\n", testWireMessage)
		app.P2p.Send(testWireMessage)
		time.Sleep(time.Minute)
	}
}

// InitServices ties the services together before running
func (app *App) InitServices(config interfaces.Config, log interfaces.Logger) {
	app.config = config
	app.log = log

	app.log.Infof("Saving data to %s", app.config.GetString("database.path"))

	// Start up the database
	app.Storage = &db.Storage{}
	app.Storage.SetDbPath(app.config.GetString("database.path"))
	app.Storage.Run()

	privateKey, publicKey, err := identity.GetIdentity(app.Storage)

	if err != nil {
		app.log.Error(err)
	}

	// Run the P2P process
	app.P2p = p2p.NewP2p(log, privateKey, publicKey)

	// Construct the server struct
	app.Server = service.NewServer(log, app.Storage, app.P2p)

	// Connect the order and channel services with p2p
	app.P2p.RegisterOrderService(app.Server.Orders)
	app.P2p.RegisterChannelService(app.Server.Channels)

	// Run the P2p service before running the gRPC server
	app.P2p.Run()
}

// Run is a separated main-function to ease testing
func (app *App) Run() {
	defer app.Storage.Close()
	defer app.P2p.Close()

	if app.config.GetBool("p2p.debug") == true {
		app.log.Info("Running the debug pinger on channel \"testChannel\"!")
		go app.debugPinger()
	}

	// Run the gRPC API
	app.Server.Run(app.config.GetUint("api.port"))
}

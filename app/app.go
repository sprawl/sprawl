package app

import (
	"time"

	config "github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/db"
	"github.com/eqlabs/sprawl/identity"
	"github.com/eqlabs/sprawl/p2p"
	"github.com/eqlabs/sprawl/pb"
	"github.com/eqlabs/sprawl/service"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
)

// App ties Sprawl's services together
type App struct {
	Storage *db.Storage
	P2p     *p2p.P2p
	Server  *service.Server
}

var appConfig *config.Config
var logger *zap.Logger
var log *zap.SugaredLogger

func init() {
	appConfig = &config.Config{}
	appConfig.ReadConfig("../config/default")
	logger, _ = zap.NewProduction()
	log = logger.Sugar()
}

func debugPinger(p2pInstance *p2p.P2p) {
	var testChannel *pb.Channel = &pb.Channel{Id: []byte("testChannel")}
	p2pInstance.Subscribe(testChannel)

	var testOrder *pb.Order = &pb.Order{Asset: string("ETH"), CounterAsset: string("BTC"), Amount: 52152, Price: 0.2, Id: []byte("jgkahgkjal")}
	testOrderInBytes, err := proto.Marshal(testOrder)
	if err != nil {
		panic(err)
	}

	testWireMessage := &pb.WireMessage{ChannelID: testChannel.GetId(), Operation: pb.Operation_CREATE, Data: testOrderInBytes}

	for {
		log.Infof("Debug pinger is sending testWireMessage: %s\n", testWireMessage)
		p2pInstance.Send(testWireMessage)
		time.Sleep(time.Minute)
	}
}

// InitServices ties the services together before running
func (app *App) InitServices() {
	log.Infof("Saving data to %s", appConfig.GetString("database.path"))

	// Start up the database
	app.Storage = &db.Storage{}
	app.Storage.SetDbPath(appConfig.GetString("database.path"))
	app.Storage.Run()

	privateKey, publicKey, err := identity.GetIdentity(app.Storage)

	if err != nil {
		log.Error(err)
	}

	// Run the P2P process
	app.P2p = p2p.NewP2p(privateKey, publicKey)

	// Construct the server struct
	app.Server = service.NewServer(app.Storage, app.P2p)

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

	if appConfig.GetBool("p2p.debug") == true {
		log.Info("Running the debug pinger on channel \"testChannel\"!")
		go debugPinger(app.P2p)
	}

	// Run the gRPC API
	app.Server.Run(appConfig.GetUint("api.port"))
}

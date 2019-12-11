package app

import (
	"context"
	"os"
	"testing"

	"github.com/sprawl/sprawl/config"
	"github.com/sprawl/sprawl/pb"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const asset1 string = "ETH"
const asset2 string = "BTC"
const testAmount = 52617562718
const testPrice = 0.1
const p2pDebugEnvVar string = "SPRAWL_P2P_DEBUG"
const envTestP2PDebug string = "true"
const testConfigPath = "../config/test"

var appConfig *config.Config
var logger *zap.Logger
var log *zap.SugaredLogger

func init() {
	logger = zap.NewNop()
	log = logger.Sugar()
	appConfig = &config.Config{}
	appConfig.ReadConfig(testConfigPath)
}

func TestApp(t *testing.T) {
	app := &App{}
	app.InitServices(appConfig, log)

	assert.NotNil(t, app.Storage)

	assert.NotNil(t, app.Server)
	assert.NotNil(t, app.Server.Orders)
	assert.NotNil(t, app.Server.Channels)

	assert.NotNil(t, app.P2p)
	assert.NotNil(t, app.P2p.Receiver)

	assert.Equal(t, app.Server.Orders, app.P2p.Receiver)

	err := app.Server.Channels.Storage.Put([]byte(asset1), []byte(asset2))
	assert.NoError(t, err)

	err = app.Server.Orders.Storage.Put([]byte(asset1), []byte(asset2))
	assert.NoError(t, err)

	ctx := context.Background()
	joinres, _ := app.Server.Channels.Join(ctx, &pb.JoinRequest{Asset: asset1, CounterAsset: asset2})
	channel := joinres.GetJoinedChannel()

	testOrder := pb.CreateRequest{ChannelID: channel.GetId(), Asset: asset1, CounterAsset: asset2, Amount: testAmount, Price: testPrice}

	_, err = app.Server.Orders.Create(ctx, &testOrder)
	assert.NoError(t, err)
	go app.Run()

	app.Storage.DeleteAll()

	defer app.Server.Close()
	defer app.Storage.Close()
	defer app.P2p.Close()
}

// TODO: doesn't test now that the debugPinger actually joins any channel. Needs refactoring of the debugPinger functionality itself to make it more testable.
func TestDebugPinger(t *testing.T) {
	app := &App{}
	os.Setenv(p2pDebugEnvVar, envTestP2PDebug)
	app.InitServices(appConfig, log)

	go app.debugPinger()

	defer app.Storage.Close()
	defer app.P2p.Close()

	os.Clearenv()
	app.Storage.DeleteAll()
}

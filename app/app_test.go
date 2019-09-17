package app

import (
	"context"
	"testing"

	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/pb"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const asset1 string = "ETH"
const asset2 string = "BTC"
const testAmount = 52617562718
const testPrice = 0.1

var appConfig *config.Config
var logger *zap.Logger
var log *zap.SugaredLogger

func init() {
	logger, _ = zap.NewProduction()
	log = logger.Sugar()
	appConfig = &config.Config{Logger: log}
	appConfig.ReadConfig("../config/default")
}

func TestApp(t *testing.T) {
	app := &App{}
	app.InitServices(appConfig, log)

	assert.NotNil(t, app.Storage)

	assert.NotNil(t, app.Server)
	assert.NotNil(t, app.Server.Orders)
	assert.NotNil(t, app.Server.Channels)

	assert.NotNil(t, app.P2p)
	assert.NotNil(t, app.P2p.Orders)
	assert.NotNil(t, app.P2p.Channels)

	assert.Equal(t, app.Server.Orders, app.P2p.Orders)
	assert.Equal(t, app.Server.Channels, app.P2p.Channels)

	err := app.Server.Channels.Storage.Put([]byte(asset1), []byte(asset2))
	assert.NoError(t, err)

	err = app.Server.Orders.Storage.Put([]byte(asset1), []byte(asset2))
	assert.NoError(t, err)

	ctx := context.Background()
	joinres, _ := app.P2p.Channels.Join(ctx, &pb.JoinRequest{Asset: asset1, CounterAsset: asset2})
	channel := joinres.GetJoinedChannel()

	testOrder := pb.CreateRequest{ChannelID: channel.GetId(), Asset: asset1, CounterAsset: asset2, Amount: testAmount, Price: testPrice}

	_, err = app.P2p.Orders.Create(ctx, &testOrder)
	assert.NoError(t, err)
}

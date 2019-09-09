package app

import (
	"context"
	"testing"

	"github.com/eqlabs/sprawl/pb"
	"github.com/stretchr/testify/assert"
)

const asset1 string = "ETH"
const asset2 string = "BTC"
const testAmount = 52617562718
const testPrice = 0.1

func TestApp(t *testing.T) {
	app := &App{}
	app.InitServices()

	assert.NotEqual(t, app.Storage, nil)

	assert.NotEqual(t, app.Server, nil)
	assert.NotEqual(t, app.Server.Orders, nil)
	assert.NotEqual(t, app.Server.Channels, nil)

	assert.NotEqual(t, app.P2p, nil)
	assert.NotEqual(t, app.P2p.Orders, nil)
	assert.NotEqual(t, app.P2p.Channels, nil)

	assert.Equal(t, app.Server.Orders, app.P2p.Orders)
	assert.Equal(t, app.Server.Channels, app.P2p.Channels)

	err := app.Server.Channels.Storage.Put([]byte(asset1), []byte(asset2))
	assert.Equal(t, err, nil)

	err = app.Server.Orders.Storage.Put([]byte(asset1), []byte(asset2))
	assert.Equal(t, err, nil)

	ctx := context.Background()
	joinres, _ := app.P2p.Channels.Join(ctx, &pb.ChannelOptions{Asset: asset1, CounterAsset: asset2})
	channel := joinres.GetJoinedChannel()

	testOrder := pb.CreateRequest{ChannelID: channel.GetId(), Asset: asset1, CounterAsset: asset2, Amount: testAmount, Price: testPrice}

	_, err = app.P2p.Orders.Create(ctx, &testOrder)
	assert.Equal(t, nil, err)
}

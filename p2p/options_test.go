package p2p

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	config "github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/identity"
	libp2p "github.com/libp2p/go-libp2p"
	libp2pConfig "github.com/libp2p/go-libp2p/config"
	"github.com/stretchr/testify/assert"
)

const testConfigPath = "../config/test"
const optionsEnableDHT string = "SPRAWL_OPTIONS_ENABLEDHT"
const optionsEnableIdentity string = "SPRAWL_OPTIONS_ENABLEIDENTITY"
const optionsEnableRelay string = "SPRAWL_OPTIONS_ENABLERELAY"
const optionsEnableAutoRelay string = "SPRAWL_OPTIONS_ENABLEAUTORELAY"
const optionsEnableNATPortMap string = "SPRAWL_OPTIONS_ENABLENATPORTMAP"

func readTestConfig() {
	// Load config
	appConfig = &config.Config{}
	appConfig.ReadConfig("../config/test")
}

func resetOptions() {
	os.Unsetenv(optionsEnableDHT)
	os.Unsetenv(optionsEnableIdentity)
	os.Unsetenv(optionsEnableRelay)
	os.Unsetenv(optionsEnableAutoRelay)
	os.Unsetenv(optionsEnableNATPortMap)
}

func TestCreateOptions(t *testing.T) {
	readTestConfig()
	privateKey, publicKey, err := identity.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err)
	p2pInstance := NewP2p(privateKey, publicKey)
	p2pInstance.initContext()
	configOptions := p2pInstance.CreateOptions()
	options := []libp2pConfig.Option{}
	options = append(options, p2pInstance.initDHT())
	options = append(options, libp2p.Identity(p2pInstance.privateKey))
	options = append(options, libp2p.EnableRelay())
	options = append(options, libp2p.EnableAutoRelay())
	options = append(options, libp2p.NATPortMap())
	assert.Equal(t, fmt.Sprintf("%v", configOptions), fmt.Sprintf("%v", options))
	os.Setenv(optionsEnableDHT, "false")
	os.Setenv(optionsEnableIdentity, "false")
	os.Setenv(optionsEnableRelay, "false")
	os.Setenv(optionsEnableAutoRelay, "false")
	os.Setenv(optionsEnableNATPortMap, "false")
	configOptions = p2pInstance.CreateOptions()
	options = []libp2pConfig.Option{}
	assert.Equal(t, fmt.Sprintf("%v", configOptions), fmt.Sprintf("%v", options))
}

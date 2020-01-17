package p2p

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	libp2p "github.com/libp2p/go-libp2p"
	libp2pConfig "github.com/libp2p/go-libp2p/config"
	config "github.com/sprawl/sprawl/config"
	"github.com/sprawl/sprawl/identity"
	"github.com/stretchr/testify/assert"
)

const testConfigPath = "../config/test"
const optionsEnableRelay string = "SPRAWL_P2P_ENABLERELAY"
const optionsEnableAutoRelay string = "SPRAWL_P2P_ENABLEAUTORELAY"
const optionsEnableNATPortMap string = "SPRAWL_P2P_ENABLENATPORTMAP"
const optionsExternalIP string = "SPRAWL_P2P_EXTERNALIP"

var appConfig *config.Config

func readTestConfig() {
	// Load config
	appConfig = &config.Config{}
	appConfig.ReadConfig(testConfigPath)
}

func resetOptions() {
	os.Unsetenv(optionsEnableRelay)
	os.Unsetenv(optionsEnableAutoRelay)
	os.Unsetenv(optionsEnableNATPortMap)
	os.Unsetenv(optionsExternalIP)
}

func TestCreateOptions(t *testing.T) {
	readTestConfig()

	privateKey, publicKey, err := identity.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err)

	p2pInstance := NewP2p(appConfig, privateKey, publicKey, Logger(log))
	p2pInstance.InitContext()

	configOptions := p2pInstance.CreateOptions()
	options := []libp2pConfig.Option{}

	options = append(options, p2pInstance.initDHT())
	options = append(options, libp2p.Identity(p2pInstance.privateKey))
	options = append(options, libp2p.EnableRelay())
	options = append(options, libp2p.EnableAutoRelay())
	options = append(options, libp2p.NATPortMap())

	assert.Equal(t, fmt.Sprintf("%v", configOptions), fmt.Sprintf("%v", options))
	resetOptions()
}

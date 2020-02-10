package p2p

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	libp2p "github.com/libp2p/go-libp2p"
	libp2pConfig "github.com/libp2p/go-libp2p/config"
	"github.com/multiformats/go-multiaddr"
	ma "github.com/multiformats/go-multiaddr"
	config "github.com/sprawl/sprawl/config"
	"github.com/sprawl/sprawl/database/inmemory"
	"github.com/sprawl/sprawl/identity"
	"github.com/sprawl/sprawl/util"
	"github.com/stretchr/testify/assert"
)

const testConfigPath = "../config/test"
const optionsEnableRelay string = "SPRAWL_P2P_ENABLERELAY"
const optionsEnableAutoRelay string = "SPRAWL_P2P_ENABLEAUTORELAY"
const optionsEnableNATPortMap string = "SPRAWL_P2P_ENABLENATPORTMAP"
const optionsExternalIP string = "SPRAWL_P2P_EXTERNALIP"
const optionsP2PPort string = "SPRAWL_P2P_PORT"

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

	testLogger := new(util.TestLogger)
	testLogger.Test(t)
	p2pInstance := &P2p{Logger: testLogger}
	p2pInstance.initDHT()
	falseOptions := []libp2pConfig.Option{}
	assert.Panics(t, func() { p2pInstance.InitHost(falseOptions...) })

	storage := &inmemory.Storage{
		Db: make(map[string]string),
	}

	p2pInstance = NewP2p(appConfig, privateKey, publicKey, Logger(testLogger), Storage(storage))

	externalIP := "127.0.0.1"
	p2pPort := "4001"

	addr, err := createMultiAddr(externalIP, p2pPort)
	assert.NoError(t, err)
	multiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf(addrTemplate, externalIP, p2pPort))
	assert.Equal(t, multiAddr, addr)

	configOptions := p2pInstance.CreateOptions()
	options := []libp2pConfig.Option{}

	options = append(options, p2pInstance.initDHT())
	options = append(options, libp2p.Identity(p2pInstance.privateKey))
	options = append(options, libp2p.EnableRelay())
	options = append(options, libp2p.EnableAutoRelay())
	options = append(options, libp2p.NATPortMap())

	assert.Equal(t, fmt.Sprintf("%v", configOptions), fmt.Sprintf("%v", options))

	os.Setenv(optionsEnableNATPortMap, "false")
	os.Setenv(optionsExternalIP, externalIP)
	os.Setenv(optionsP2PPort, p2pPort)
	appConfig.ReadConfig(testConfigPath)
	customIPOptions := p2pInstance.CreateOptions()

	options = []libp2pConfig.Option{}

	options = append(options, p2pInstance.initDHT())
	options = append(options, libp2p.Identity(p2pInstance.privateKey))
	options = append(options, libp2p.EnableRelay())
	options = append(options, libp2p.EnableAutoRelay())
	multiaddrs := []ma.Multiaddr{}
	multiaddrs = append(multiaddrs, multiAddr)
	addrFactory := func(addrs []ma.Multiaddr) []ma.Multiaddr {
		return multiaddrs
	}
	options = append(options, libp2p.ListenAddrs(multiaddrs...))
	options = append(options, libp2p.AddrsFactory(addrFactory))

	assert.Equal(t, fmt.Sprintf("%v", customIPOptions), fmt.Sprintf("%v", options))

	resetOptions()
}

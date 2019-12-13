package p2p

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	libp2p "github.com/libp2p/go-libp2p"
	libp2pConfig "github.com/libp2p/go-libp2p/config"
	ma "github.com/multiformats/go-multiaddr"
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
	p2pInstance.initContext()

	configOptions := p2pInstance.CreateOptions()
	options := []libp2pConfig.Option{}

	options = append(options, p2pInstance.initDHT())
	options = append(options, libp2p.Identity(p2pInstance.privateKey))
	options = append(options, libp2p.EnableRelay())
	options = append(options, libp2p.EnableAutoRelay())
	multiaddrs := defaultListenAddrs(appConfig.GetP2PPort())
	addrFactory := func(addrs []ma.Multiaddr) []ma.Multiaddr {
		return multiaddrs
	}
	options = append(options, libp2p.ListenAddrs(multiaddrs...))
	options = append(options, libp2p.AddrsFactory(addrFactory))
	assert.Equal(t, fmt.Sprintf("%v", configOptions), fmt.Sprintf("%v", options))

	options = options[:len(options)-2]
	externalIP := "192.168.0.1"
	os.Setenv(optionsExternalIP, externalIP)
	multiaddrs = defaultListenAddrs(appConfig.GetP2PPort())
	externalMultiaddr, err := createMultiAddr(externalIP, appConfig.GetP2PPort())
	assert.Nil(t, err)
	multiaddrs = append(multiaddrs, externalMultiaddr)
	addrFactory = func(addrs []ma.Multiaddr) []ma.Multiaddr {
		return multiaddrs
	}
	options = append(options, libp2p.ListenAddrs(multiaddrs...))
	options = append(options, libp2p.AddrsFactory(addrFactory))
	configOptions = p2pInstance.CreateOptions()
	assert.Equal(t, fmt.Sprintf("%v", configOptions), fmt.Sprintf("%v", options))

	os.Setenv(optionsEnableNATPortMap, "true")
	configOptions = p2pInstance.CreateOptions()
	options = options[:len(options)-2]
	options = append(options, libp2p.NATPortMap())
	assert.Equal(t, fmt.Sprintf("%v", configOptions), fmt.Sprintf("%v", options))

	resetOptions()
}

package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ListenPort             string `toml:"listen_port"`
	ListenAddr             string `toml:"listen_addr"`
	Bootstrap              bool   `toml:"is_bootstrap"`
	BootstrapNodeMultiaddr string `toml:"bootstrap_node_multiaddr"`
	Rendezvous             string `toml:"rendezvous"`
	ProtocolID             string `toml:"protocol_id"`
	SessionKey             string `toml:"session_key"`
	PrivateKey             string `toml:"private_key"`
}

// viperEnvVariable loads config parameters from .env file
func viperEnvString(key string, default_value string) string {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		return default_value
	}

	return value
}

func viperEnvBoolean(key string, default_value bool) bool {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	value := viper.GetBool(key)

	return value
}

// NewConfig creates a new config based on default values or provided .env file
func NewConfig() *Config {
	ListenPort := viperEnvString("LISTEN_PORT", ":8000")
	ListenAddr := viperEnvString("LISTEN_ADDRESS", "debug")
	Bootstrap := viperEnvBoolean("BOOTSTRAP_NODE", false)
	BootstrapNodeMultiaddr := viperEnvString("BOOTSTRAP_NODE_MULTIADDRESS", "/ip4/127.0.0.1/tcp/0")
	Rendezvous := viperEnvString("RENDEZVOUS", "filecoin-p2p-oracle")
	ProtocolID := viperEnvString("PROTOCOL_ID", "p2p-oracle")
	SessionKey := viperEnvString("SESSION_KEY", "go")
	PrivateKey := viperEnvString("PRIVATE_KEY", "")

	return &Config{
		ListenPort:             ListenPort,
		ListenAddr:             ListenAddr,
		Bootstrap:              Bootstrap,
		BootstrapNodeMultiaddr: BootstrapNodeMultiaddr,
		Rendezvous:             Rendezvous,
		ProtocolID:             ProtocolID,
		SessionKey:             SessionKey,
		PrivateKey:             PrivateKey,
	}
}

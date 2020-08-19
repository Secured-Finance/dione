package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ListenPort             string         `mapstructure:"listen_port"`
	ListenAddr             string         `mapstructure:"listen_addr"`
	Bootstrap              bool           `mapstructure:"is_bootstrap"`
	BootstrapNodeMultiaddr string         `mapstructure:"bootstrap_node_multiaddr"`
	Rendezvous             string         `mapstructure:"rendezvous"`
	SessionKey             string         `mapstructure:"session_key"`
	Etherium               EtheriumConfig `mapstructure:"eth"`
	PubSub                 PubSubConfig   `mapstructure:"pubSub"`
}

type EtheriumConfig struct {
	PrivateKey string `mapstructure:"private_key"`
}

type PubSubConfig struct {
	ProtocolID string `mapstructure:"protocolID"`
}

// NewConfig creates a new config based on default values or provided .env file
func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{
		ListenAddr:             "localhost",
		ListenPort:             ":8000",
		Bootstrap:              false,
		BootstrapNodeMultiaddr: "/ip4/127.0.0.1/tcp/0",
		Rendezvous:             "filecoin-p2p-oracle",
		Etherium: EtheriumConfig{
			PrivateKey: "",
		},
		SessionKey: "go",
		PubSub: PubSubConfig{
			ProtocolID: "p2p-oracle",
		},
	}

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

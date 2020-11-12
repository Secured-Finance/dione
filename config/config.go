package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ListenPort             int            `mapstructure:"listen_port"`
	ListenAddr             string         `mapstructure:"listen_addr"`
	BootstrapNodes         []string       `mapstructure:"bootstrap_node_multiaddr"`
	Rendezvous             string         `mapstructure:"rendezvous"`
	Ethereum               EthereumConfig `mapstructure:"ethereum"`
	Filecoin               FilecoinConfig `mapstructure:"filecoin"`
	PubSub                 PubSubConfig   `mapstructure:"pubSub"`
	Store                  StoreConfig    `mapstructure:"store"`
	ConsensusMaxFaultNodes int            `mapstructure:"consensus_max_fault_nodes"`
}

type EthereumConfig struct {
	GatewayAddress               string `mapstructure:"gateway_address"`
	PrivateKey                   string `mapstructure:"private_key"`
	OracleEmitterContractAddress string `mapstructure:"oracle_emitter_contract_address"`
	AggregatorContractAddress    string `mapstructure:"aggregator_contract_address"`
}

type FilecoinConfig struct {
	LotusHost  string `mapstructure:"lotusHost"`
	LotusToken string `mapstructure:"lotusToken"`
}

type PubSubConfig struct {
	ProtocolID string `mapstructure:"protocolID"`
}

type StoreConfig struct {
	DatabaseURL string `mapstructure:"database_url"`
}

// NewConfig creates a new config based on default values or provided .env file
func NewConfig(configPath string) (*Config, error) {
	dbName := "dione"
	username := "user"
	password := "password"
	dbURL := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable", username, password, dbName)

	cfg := &Config{
		ListenAddr:     "localhost",
		ListenPort:     8000,
		BootstrapNodes: []string{"/ip4/127.0.0.1/tcp/0"},
		Rendezvous:     "filecoin-p2p-oracle",
		Ethereum: EthereumConfig{
			PrivateKey: "",
		},
		PubSub: PubSubConfig{
			ProtocolID: "p2p-oracle",
		},
		Store: StoreConfig{
			DatabaseURL: dbURL,
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

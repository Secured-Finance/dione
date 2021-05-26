package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ListenPort            int              `mapstructure:"listen_port"`
	ListenAddr            string           `mapstructure:"listen_addr"`
	IsBootstrap           bool             `mapstructure:"is_bootstrap"`
	BootstrapNodes        []string         `mapstructure:"bootstrap_node_multiaddr"`
	Rendezvous            string           `mapstructure:"rendezvous"`
	Ethereum              EthereumConfig   `mapstructure:"ethereum"`
	Filecoin              FilecoinConfig   `mapstructure:"filecoin"`
	PubSub                PubSubConfig     `mapstructure:"pubsub"`
	Store                 StoreConfig      `mapstructure:"store"`
	ConsensusMinApprovals int              `mapstructure:"consensus_min_approvals"`
	Redis                 RedisConfig      `mapstructure:"redis"`
	CacheType             string           `mapstructure:"cache_type"`
	Blockchain            BlockchainConfig `mapstructure:"blockchain"`
}

type EthereumConfig struct {
	GatewayAddress              string `mapstructure:"gateway_address"`
	ChainID                     int    `mapstructure:"chain_id"`
	PrivateKey                  string `mapstructure:"private_key"`
	MnemonicPhrase              string `mapstructure:"mnemonic_phrase"`
	HDDerivationPath            string `mapstructure:"hd_derivation_path"`
	DioneOracleContractAddress  string `mapstructure:"oracle_contract_address"`
	DioneStakingContractAddress string `mapstructure:"staking_contract_address"`
	DisputeContractAddress      string `mapstructure:"dispute_contract_address"`
	DisputeVoteWindow           int    `mapstructure:"dispute_vote_window"` // in secs
}

type FilecoinConfig struct {
	LotusHost  string `mapstructure:"lotusHost"`
	LotusToken string `mapstructure:"lotusToken"`
}

type PubSubConfig struct {
	ServiceTopicName string `mapstructure:"service_topic_name"`
}

type StoreConfig struct {
	DatabaseURL string `mapstructure:"database_url"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"redis_addr"`
	Password string `mapstructure:"redis_password"`
	DB       int    `mapstructure:"redis_db"`
}

type BlockchainConfig struct {
	DatabasePath string `mapstructure:"database_path"`
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
		Rendezvous:     "dione",
		Ethereum: EthereumConfig{
			PrivateKey: "",
		},
		PubSub: PubSubConfig{
			ServiceTopicName: "dione",
		},
		Store: StoreConfig{
			DatabaseURL: dbURL,
		},
		Redis: RedisConfig{
			Addr:     "redisDB:6379",
			Password: "",
			DB:       0,
		},
		CacheType: "in-memory",
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

package node

type Config struct {
	EthPrivateKey    string       `toml:"ethPrivateKey"`
	ListenAddress    string       `toml:"listenAddress"`
	ListenPort       string       `toml:"listenPort"`
	BootstrapAddress string       `toml:"bootstrapAddress"`
	RendezvousString string       `toml:"rendezvousString"`
	PubSub           PubSubConfig `toml:"pubSub"`
}

type PubSubConfig struct {
	ProtocolID string `toml:"protocolID"`
}

type EthereumConfig struct {
	PrivateKey     string `toml:"privateKey"`
	GatewayAddress string `toml:"gatewayAddress"`
}

module github.com/Secured-Finance/dione

go 1.14

require (
	github.com/deckarep/golang-set v1.7.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/ethereum/go-ethereum v1.9.5
	github.com/filecoin-project/go-address v0.0.2-0.20200504173055-8b6f2fb2b3ef
	github.com/filecoin-project/lotus v0.4.2
	github.com/filecoin-project/specs-actors v0.8.1-0.20200723200253-a3c01bc62f99
	github.com/gballet/go-libpcsclite v0.0.0-20191108122812-4678299bea08 // indirect
	github.com/google/uuid v1.1.1
	github.com/hashicorp/raft v1.1.2
	github.com/ipfs/go-log v1.0.4
	github.com/karalabe/usb v0.0.0-20191104083709-911d15fe12a9 // indirect
	github.com/libp2p/go-libp2p v0.10.2
	github.com/libp2p/go-libp2p-consensus v0.0.1
	github.com/libp2p/go-libp2p-core v0.6.1
	github.com/libp2p/go-libp2p-discovery v0.5.0
	github.com/libp2p/go-libp2p-kad-dht v0.8.3
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/libp2p/go-libp2p-pubsub v0.3.3
	github.com/libp2p/go-libp2p-raft v0.1.5
	github.com/multiformats/go-multiaddr v0.2.2
	github.com/olekukonko/tablewriter v0.0.4 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/renproject/hyperdrive v1.2.0
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/spf13/viper v1.7.1
	github.com/status-im/keycard-go v0.0.0-20200402102358-957c09536969 // indirect
	github.com/tyler-smith/go-bip39 v1.0.2 // indirect
	github.com/wsddn/go-ecdh v0.0.0-20161211032359-48726bab9208 // indirect
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9
)

replace github.com/libp2p/go-libp2p-raft v0.1.5 => github.com/ItalyPaleAle/go-libp2p-raft v0.1.6-0.20200703060436-0b37aa16095e

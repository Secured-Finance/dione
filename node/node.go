package node

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"

	"github.com/Secured-Finance/p2p-oracle-node/config"
	consensus "github.com/Secured-Finance/p2p-oracle-node/consensus"
	"github.com/Secured-Finance/p2p-oracle-node/handler"
	"github.com/Secured-Finance/p2p-oracle-node/rpc"
	"github.com/Secured-Finance/p2p-oracle-node/rpcclient"
	mapset "github.com/deckarep/golang-set"
	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
)

type Node struct {
	Host            host.Host
	PubSub          *pubsub.PubSub
	GlobalCtx       context.Context
	GlobalCtxCancel context.CancelFunc
	OracleTopic     string
	networkTopics   mapset.Set
	handler         *handler.Handler
	Config          *config.Config
	Logger          *log.ZapEventLogger
	Consensus       *consensus.RaftConsensus
	Lotus           *rpc.LotusClient
	Ethereum        *rpcclient.EthereumClient
}

func NewNode(configPath string) (*Node, error) {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return nil, err
	}
	node := &Node{
		OracleTopic:   "p2p_oracle",
		Config:        cfg,
		Logger:        log.Logger("node"),
		networkTopics: mapset.NewSet(),
	}
	log.SetAllLoggers(log.LevelInfo)

	return node, nil
}

func (node *Node) setupNode(ctx context.Context, prvKey crypto.PrivKey) {
	listenMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", node.Config.ListenAddr, node.Config.ListenPort))
	if err != nil {
		node.Logger.Fatal("Failed to generate new node multiaddress:", err)
	}
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(listenMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		node.Logger.Fatal("Failed to set a new libp2p node:", err)
	}
	node.Host = host
	node.startPubSub(ctx, host)
}

func Start() error {
	configPath := flag.String("config", "", "Path to config")
	verbose := flag.Bool("verbose", false, "Verbose logging")
	flag.Parse()

	if *configPath == "" {
		return fmt.Errorf("no config path provided")
	}

	node, err := NewNode(*configPath)
	if *verbose {
		log.SetAllLoggers(log.LevelDebug)
	} else {
		log.SetAllLoggers(log.LevelInfo)
	}
	if err != nil {
		log.Logger("node").Panic(err)
	}

	r := rand.Reader

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		node.Logger.Fatal(err)
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	node.GlobalCtx = ctx
	node.GlobalCtxCancel = ctxCancel

	node.setupNode(ctx, prvKey)
	return nil
}

package node

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	pex "github.com/Secured-Finance/go-libp2p-pex"
	"time"

	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/consensus"
	"github.com/Secured-Finance/dione/pb"
	"github.com/Secured-Finance/dione/rpc"
	"github.com/Secured-Finance/dione/rpcclient"
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

const (
	DefaultPEXUpdateTime = 1*time.Minute
)

type Node struct {
	Host             host.Host
	PubSubRouter     *pb.PubSubRouter
	GlobalCtx        context.Context
	GlobalCtxCancel  context.CancelFunc
	OracleTopic      string
	Config           *config.Config
	Lotus            *rpc.LotusClient
	Ethereum         *rpcclient.EthereumClient
	ConsensusManager *consensus.PBFTConsensusManager
}

func NewNode(configPath string) (*Node, error) {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return nil, err
	}
	node := &Node{
		OracleTopic: "dione",
		Config:      cfg,
	}

	return node, nil
}

func (n *Node) setupNode(ctx context.Context, prvKey crypto.PrivKey, pexDiscoveryUpdateTime time.Duration, maxFaultNodes int) {
	n.setupLibp2pHost(context.TODO(), prvKey, pexDiscoveryUpdateTime)
	//n.setupEthereumClient()
	//n.setupFilecoinClient()
	n.setupPubsub()
	n.setupConsensusManager(maxFaultNodes)
}

func (n *Node) setupEthereumClient() error {
	ethereum := rpcclient.NewEthereumClient()
	n.Ethereum = ethereum
	return ethereum.Initialize(context.Background(),
		n.Config.Ethereum.GatewayAddress,
		n.Config.Ethereum.PrivateKey,
		n.Config.Ethereum.OracleEmitterContractAddress,
		n.Config.Ethereum.AggregatorContractAddress,
	)
}

func (n *Node) setupFilecoinClient() {
	lotus := rpc.NewLotusClient(n.Config.Filecoin.LotusHost, n.Config.Filecoin.LotusToken)
	n.Lotus = lotus
}

func (n *Node) setupPubsub() {
	n.PubSubRouter = pb.NewPubSubRouter(n.Host, n.OracleTopic)
	// wait for setting up pubsub
	//time.Sleep(3 * time.Second)
}

func (n *Node) setupConsensusManager(maxFaultNodes int) {
	n.ConsensusManager = consensus.NewPBFTConsensusManager(n.PubSubRouter, maxFaultNodes)
}

func (n *Node) setupLibp2pHost(ctx context.Context, privateKey crypto.PrivKey, pexDiscoveryUpdateTime time.Duration) {
	listenMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", n.Config.ListenAddr, n.Config.ListenPort))
	if err != nil {
		logrus.Fatal("Failed to generate new node multiaddress:", err)
	}
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(listenMultiAddr),
		libp2p.Identity(privateKey),
	)
	if err != nil {
		logrus.Fatal("Failed to set a new libp2p node:", err)
	}
	n.Host = host

	logrus.Info(fmt.Sprintf("[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s", n.Config.ListenAddr, n.Config.ListenPort, host.ID().Pretty()))

	var bootstrapMaddrs []multiaddr.Multiaddr
	for _, a := range n.Config.BootstrapNodes {
		maddr, err := multiaddr.NewMultiaddr(a)
		if err != nil {
			logrus.Fatalf("Invalid multiaddress of bootstrap node: %s", err.Error())
		}
		bootstrapMaddrs = append(bootstrapMaddrs, maddr)
	}

	discovery, err := pex.NewPEXDiscovery(host, bootstrapMaddrs, pexDiscoveryUpdateTime)

	logrus.Info("Announcing ourselves...")
	_, err = discovery.Advertise(context.TODO(), n.Config.Rendezvous)
	if err != nil {
		logrus.Fatalf("Failed to announce this node to the network: %s", err.Error())
	}
	logrus.Info("Successfully announced!")

	// Discover unbounded count of peers
	logrus.Info("Searching for other peers...")
	peerChan, err := discovery.FindPeers(context.TODO(), n.Config.Rendezvous)
	if err != nil {
		logrus.Fatal("Failed to find new peers, exiting...", err)
	}
	go func() {
	MainLoop:
		for {
			select {
			case <-ctx.Done():
				break MainLoop
			case newPeer := <-peerChan:
				{
					if len(newPeer.Addrs) == 0 {
						continue
					}
					if newPeer.ID.String() == n.Host.ID().String() {
						continue
					}
					logrus.Infof("Found peer: %s", newPeer)
					// Connect to the peer
					if err := n.Host.Connect(ctx, newPeer); err != nil {
						logrus.Warn("Connection failed: ", err)
					}
					logrus.Info("Connected to newly discovered peer: ", newPeer)
				}
			}
		}
	}()
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
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	if err != nil {
		logrus.Panic(err)
	}

	privKey, err := generatePrivateKey()
	if err != nil {
		logrus.Fatal(err)
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	node.GlobalCtx = ctx
	node.GlobalCtxCancel = ctxCancel

	node.setupNode(ctx, privKey, DefaultPEXUpdateTime, node.Config.ConsensusMaxFaultNodes)
	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func generatePrivateKey() (crypto.PrivKey, error) {
	r := rand.Reader
	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}
	return prvKey, nil
}

// TODO generate MinerBase for the node

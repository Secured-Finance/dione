package node

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/Secured-Finance/dione/rpc"
	rtypes "github.com/Secured-Finance/dione/rpc/types"

	solana2 "github.com/Secured-Finance/dione/rpc/solana"

	"github.com/Secured-Finance/dione/rpc/filecoin"

	"github.com/Secured-Finance/dione/types"

	"github.com/Secured-Finance/dione/wallet"

	"golang.org/x/xerrors"

	"github.com/Secured-Finance/dione/beacon"

	pex "github.com/Secured-Finance/go-libp2p-pex"

	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/consensus"
	"github.com/Secured-Finance/dione/ethclient"
	"github.com/Secured-Finance/dione/pubsub"
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

const (
	DefaultPEXUpdateTime = 6 * time.Second
)

type Node struct {
	Host             host.Host
	PubSubRouter     *pubsub.PubSubRouter
	GlobalCtx        context.Context
	GlobalCtxCancel  context.CancelFunc
	OracleTopic      string
	Config           *config.Config
	Ethereum         *ethclient.EthereumClient
	ConsensusManager *consensus.PBFTConsensusManager
	Miner            *consensus.Miner
	Beacon           beacon.BeaconNetworks
	Wallet           *wallet.LocalWallet
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

func (n *Node) setupNode(ctx context.Context, prvKey crypto.PrivKey, pexDiscoveryUpdateTime time.Duration) {
	n.setupLibp2pHost(context.TODO(), prvKey, pexDiscoveryUpdateTime)
	err := n.setupEthereumClient()
	if err != nil {
		logrus.Fatal(err)
	}

	err = n.setupRPCClients()
	if err != nil {
		logrus.Fatal(err)
	}

	n.setupPubsub()
	rawPrivKey, err := prvKey.Raw()
	if err != nil {
		logrus.Fatal(err)
	}
	err = n.setupBeacon()
	if err != nil {
		logrus.Fatal(err)
	}
	err = n.setupMiner(rawPrivKey)
	if err != nil {
		logrus.Fatal(err)
	}
	err = n.setupConsensusManager(rawPrivKey, n.Config.ConsensusMinApprovals)
	if err != nil {
		logrus.Fatalf("Failed to setup consensus manager: %w", err)
	}
	err = n.setupWallet(rawPrivKey)
	if err != nil {
		logrus.Fatal(err)
	}
	n.subscribeOnEthContracts(ctx)
}

func (n *Node) setupMiner(privKey []byte) error {
	n.Miner = consensus.NewMiner(n.Host.ID(), *n.Ethereum.GetEthAddress(), n.Beacon, n.Ethereum, privKey)
	return nil
}

func (n *Node) setupBeacon() error {
	beacon, err := n.NewBeaconClient()
	if err != nil {
		return xerrors.Errorf("failed to setup beacon: %w", err)
	}
	n.Beacon = beacon
	return nil
}

func (n *Node) setupWallet(privKey []byte) error {
	// TODO make persistent keystore
	kstore := wallet.NewMemKeyStore()
	keyInfo := types.KeyInfo{
		Type:       types.KTEd25519,
		PrivateKey: privKey,
	}

	kstore.Put(wallet.KNamePrefix+n.Host.ID().String(), keyInfo)
	w, err := wallet.NewWallet(kstore)
	if err != nil {
		return xerrors.Errorf("failed to setup wallet: %w", err)
	}
	n.Wallet = w
	return nil
}

func (n *Node) setupEthereumClient() error {
	ethereum := ethclient.NewEthereumClient()
	n.Ethereum = ethereum
	return ethereum.Initialize(context.Background(),
		n.Config.Ethereum.GatewayAddress,
		n.Config.Ethereum.PrivateKey,
		n.Config.Ethereum.OracleEmitterContractAddress,
		n.Config.Ethereum.AggregatorContractAddress,
		n.Config.Ethereum.DioneStakingContractAddress,
	)
}

func (n *Node) setupRPCClients() error {
	fc := filecoin.NewLotusClient()
	rpc.RegisterRPC(rtypes.RPCTypeFilecoin, map[string]func(string) (string, error){
		"getTransaction": fc.GetTransaction,
	})

	sl := solana2.NewSolanaClient()
	rpc.RegisterRPC(rtypes.RPCTypeSolana, map[string]func(string) (string, error){
		"getTransaction": sl.GetTransaction,
	})

	return nil
}

func (n *Node) setupPubsub() {
	n.PubSubRouter = pubsub.NewPubSubRouter(n.Host, n.OracleTopic)
	// wait for setting up pubsub
	//time.Sleep(3 * time.Second)
}

func (n *Node) setupConsensusManager(privateKey []byte, minApprovals int) error {
	n.ConsensusManager = consensus.NewPBFTConsensusManager(n.PubSubRouter, minApprovals, privateKey, n.Ethereum, n.Miner)
	return nil
}

func (n *Node) setupLibp2pHost(ctx context.Context, privateKey crypto.PrivKey, pexDiscoveryUpdateTime time.Duration) {
	listenMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", n.Config.ListenAddr, n.Config.ListenPort))
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

	logrus.Info(fmt.Sprintf("[*] Your Multiaddress Is: /ip4/%s/tcp/%d/p2p/%s", n.Config.ListenAddr, n.Config.ListenPort, host.ID().Pretty()))

	var bootstrapMaddrs []multiaddr.Multiaddr
	for _, a := range n.Config.BootstrapNodes {
		maddr, err := multiaddr.NewMultiaddr(a)
		if err != nil {
			logrus.Fatalf("Invalid multiaddress of bootstrap node: %s", err.Error())
		}
		bootstrapMaddrs = append(bootstrapMaddrs, maddr)
	}
	if n.Config.IsBootstrap {
		bootstrapMaddrs = nil
	}
	discovery, err := pex.NewPEXDiscovery(host, bootstrapMaddrs, pexDiscoveryUpdateTime)
	if err != nil {
		logrus.Fatal("Can't set up PEX discovery protocol, exiting... ", err)
	}

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
		logrus.SetLevel(logrus.DebugLevel)
	}
	if err != nil {
		logrus.Panic(err)
	}

	var privateKey crypto.PrivKey

	if node.Config.IsBootstrap {
		if _, err := os.Stat(".bootstrap_privkey"); os.IsNotExist(err) {
			privateKey, err = generatePrivateKey()
			if err != nil {
				logrus.Fatal(err)
			}

			f, _ := os.Create(".bootstrap_privkey")
			r, _ := privateKey.Raw()
			f.Write(r)
		} else {
			pkey, _ := ioutil.ReadFile(".bootstrap_privkey")
			privateKey, _ = crypto.UnmarshalEd25519PrivateKey(pkey)
		}
	} else {
		privateKey, err = generatePrivateKey()
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	node.GlobalCtx = ctx
	node.GlobalCtxCancel = ctxCancel

	node.setupNode(ctx, privateKey, DefaultPEXUpdateTime)
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
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, r)
	if err != nil {
		return nil, err
	}
	return prvKey, nil
}

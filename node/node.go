package node

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	pex "github.com/Secured-Finance/go-libp2p-pex"

	"github.com/Secured-Finance/dione/cache"

	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/Secured-Finance/dione/drand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/libp2p/go-libp2p-core/discovery"

	"github.com/Secured-Finance/dione/rpc"
	rtypes "github.com/Secured-Finance/dione/rpc/types"

	solana2 "github.com/Secured-Finance/dione/rpc/solana"

	"github.com/Secured-Finance/dione/rpc/filecoin"

	"github.com/Secured-Finance/dione/types"

	"github.com/Secured-Finance/dione/wallet"

	"golang.org/x/xerrors"

	"github.com/Secured-Finance/dione/beacon"

	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/consensus"
	"github.com/Secured-Finance/dione/ethclient"
	pubsub2 "github.com/Secured-Finance/dione/pubsub"
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
	PeerDiscovery    discovery.Discovery
	PubSubRouter     *pubsub2.PubSubRouter
	GlobalCtx        context.Context
	GlobalCtxCancel  context.CancelFunc
	Config           *config.Config
	Ethereum         *ethclient.EthereumClient
	ConsensusManager *consensus.PBFTConsensusManager
	Miner            *consensus.Miner
	Beacon           beacon.BeaconNetworks
	Wallet           *wallet.LocalWallet
	EventLogCache    *cache.EventLogCache
}

func NewNode(config *config.Config, prvKey crypto.PrivKey, pexDiscoveryUpdateTime time.Duration) (*Node, error) {
	n := &Node{
		Config: config,
	}

	// initialize libp2p host
	lhost, err := provideLibp2pHost(n.Config, prvKey, pexDiscoveryUpdateTime)
	if err != nil {
		logrus.Fatal(err)
	}
	n.Host = lhost

	// initialize ethereum client
	ethClient, err := provideEthereumClient(n.Config)
	if err != nil {
		logrus.Fatal(err)
	}
	n.Ethereum = ethClient

	// initialize blockchain rpc clients
	err = n.setupRPCClients()
	if err != nil {
		logrus.Fatal(err)
	}

	// initialize pubsub subsystem
	psb := providePubsubRouter(lhost, n.Config)
	n.PubSubRouter = psb

	// initialize peer discovery
	peerDiscovery, err := providePeerDiscovery(n.Config, lhost, pexDiscoveryUpdateTime)
	if err != nil {
		logrus.Fatal(err)
	}
	n.PeerDiscovery = peerDiscovery

	// get private key of libp2p host
	rawPrivKey, err := prvKey.Raw()
	if err != nil {
		logrus.Fatal(err)
	}

	// initialize random beacon network subsystem
	randomBeaconNetwork, err := provideBeacon(psb.Pubsub)
	if err != nil {
		logrus.Fatal(err)
	}
	n.Beacon = randomBeaconNetwork

	// initialize mining subsystem
	miner := provideMiner(n.Host.ID(), *n.Ethereum.GetEthAddress(), n.Beacon, n.Ethereum, rawPrivKey)
	n.Miner = miner

	// initialize event log cache subsystem
	eventLogCache := provideEventLogCache()
	n.EventLogCache = eventLogCache

	// initialize consensus subsystem
	cManager := provideConsensusManager(psb, miner, ethClient, rawPrivKey, n.Config.ConsensusMinApprovals, eventLogCache)
	n.ConsensusManager = cManager

	// initialize internal eth wallet
	wallet, err := provideWallet(n.Host.ID(), rawPrivKey)
	if err != nil {
		logrus.Fatal(err)
	}
	n.Wallet = wallet

	return n, nil
}

func (n *Node) Run(ctx context.Context) error {
	n.runLibp2pAsync(ctx)
	n.subscribeOnEthContractsAsync(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}

	// return nil
}

func (n *Node) runLibp2pAsync(ctx context.Context) error {
	logrus.Info(fmt.Sprintf("[*] Your Multiaddress Is: /ip4/%s/tcp/%d/p2p/%s", n.Config.ListenAddr, n.Config.ListenPort, n.Host.ID().Pretty()))

	logrus.Info("Announcing ourselves...")
	_, err := n.PeerDiscovery.Advertise(context.TODO(), n.Config.Rendezvous)
	if err != nil {
		return xerrors.Errorf("failed to announce this node to the network: %v", err)
	}
	logrus.Info("Successfully announced!")

	// Discover unbounded count of peers
	logrus.Info("Searching for other peers...")
	peerChan, err := n.PeerDiscovery.FindPeers(context.TODO(), n.Config.Rendezvous)
	if err != nil {
		return xerrors.Errorf("failed to find new peers: %v", err)
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
	return nil
}

func (n *Node) subscribeOnEthContractsAsync(ctx context.Context) {
	eventChan, subscription, err := n.Ethereum.SubscribeOnOracleEvents(ctx)
	if err != nil {
		logrus.Fatal("Couldn't subscribe on ethereum contracts, exiting... ", err)
	}

	go func() {
	EventLoop:
		for {
			select {
			case event := <-eventChan:
				{
					err := n.EventLogCache.Store("request_"+event.RequestID.String(), event)
					if err != nil {
						logrus.Errorf("Failed to store new request event to event log cache: %v", err)
					}

					logrus.Info("Let's wait a little so that all nodes have time to receive the request and cache it")
					time.Sleep(5 * time.Second)

					task, err := n.Miner.MineTask(context.TODO(), event)
					if err != nil {
						logrus.Fatal("Failed to mine task, exiting... ", err)
					}
					if task == nil {
						continue
					}
					logrus.Infof("Proposed new Dione task with ID: %s", event.RequestID.String())
					err = n.ConsensusManager.Propose(event.RequestID.String(), *task, event)
					if err != nil {
						logrus.Errorf("Failed to propose task: %w", err)
					}
				}
			case <-ctx.Done():
				break EventLoop
			case <-subscription.Err():
				logrus.Fatal("Error with ethereum subscription, exiting... ", err)
			}
		}
	}()
}

func provideEventLogCache() *cache.EventLogCache {
	return cache.NewEventLogCache()
}

func provideMiner(peerID peer.ID, ethAddress common.Address, beacon beacon.BeaconNetworks, ethClient *ethclient.EthereumClient, privateKey []byte) *consensus.Miner {
	return consensus.NewMiner(peerID, ethAddress, beacon, ethClient, privateKey)
}

func provideBeacon(ps *pubsub.PubSub) (beacon.BeaconNetworks, error) {
	networks := beacon.BeaconNetworks{}
	bc, err := drand.NewDrandBeacon(config.ChainGenesis, config.TaskEpochInterval, ps)
	if err != nil {
		return nil, fmt.Errorf("failed to setup drand beacon: %w", err)
	}
	networks = append(networks, beacon.BeaconNetwork{Start: types.DrandRound(config.ChainGenesis), Beacon: bc})
	// NOTE: currently we use only one network
	return networks, nil
}

func provideWallet(peerID peer.ID, privKey []byte) (*wallet.LocalWallet, error) {
	// TODO make persistent keystore
	kstore := wallet.NewMemKeyStore()
	keyInfo := types.KeyInfo{
		Type:       types.KTEd25519,
		PrivateKey: privKey,
	}

	kstore.Put(wallet.KNamePrefix+peerID.String(), keyInfo)
	w, err := wallet.NewWallet(kstore)
	if err != nil {
		return nil, xerrors.Errorf("failed to setup wallet: %w", err)
	}
	return w, nil
}

func provideEthereumClient(config *config.Config) (*ethclient.EthereumClient, error) {
	ethereum := ethclient.NewEthereumClient()
	err := ethereum.Initialize(context.Background(),
		config.Ethereum.GatewayAddress,
		config.Ethereum.PrivateKey,
		config.Ethereum.OracleEmitterContractAddress,
		config.Ethereum.AggregatorContractAddress,
		config.Ethereum.DioneStakingContractAddress,
	)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize ethereum client: %v", err)
	}
	return ethereum, nil
}

func (n *Node) setupRPCClients() error {
	fc := filecoin.NewLotusClient()
	rpc.RegisterRPC(rtypes.RPCTypeFilecoin, map[string]func(string) ([]byte, error){
		"getTransaction": fc.GetTransaction,
	})

	sl := solana2.NewSolanaClient()
	rpc.RegisterRPC(rtypes.RPCTypeSolana, map[string]func(string) ([]byte, error){
		"getTransaction": sl.GetTransaction,
	})

	return nil
}

func providePubsubRouter(lhost host.Host, config *config.Config) *pubsub2.PubSubRouter {
	return pubsub2.NewPubSubRouter(lhost, config.PubSub.ServiceTopicName, config.IsBootstrap)
}

func provideConsensusManager(psb *pubsub2.PubSubRouter, miner *consensus.Miner, ethClient *ethclient.EthereumClient, privateKey []byte, minApprovals int, evc *cache.EventLogCache) *consensus.PBFTConsensusManager {
	return consensus.NewPBFTConsensusManager(psb, minApprovals, privateKey, ethClient, miner, evc)
}

func provideLibp2pHost(config *config.Config, privateKey crypto.PrivKey, pexDiscoveryUpdateTime time.Duration) (host.Host, error) {
	listenMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", config.ListenAddr, config.ListenPort))
	if err != nil {
		return nil, xerrors.Errorf("failed to parse multiaddress: %v", err)
	}
	host, err := libp2p.New(
		context.TODO(),
		libp2p.ListenAddrs(listenMultiAddr),
		libp2p.Identity(privateKey),
	)
	if err != nil {
		return nil, xerrors.Errorf("failed to setup libp2p host: %v", err)
	}

	return host, nil
}

func providePeerDiscovery(config *config.Config, h host.Host, pexDiscoveryUpdateTime time.Duration) (discovery.Discovery, error) {
	var bootstrapMaddrs []multiaddr.Multiaddr
	for _, a := range config.BootstrapNodes {
		maddr, err := multiaddr.NewMultiaddr(a)
		if err != nil {
			return nil, xerrors.Errorf("invalid multiaddress of bootstrap node: %v", err)
		}
		bootstrapMaddrs = append(bootstrapMaddrs, maddr)
	}

	if config.IsBootstrap {
		bootstrapMaddrs = nil
	}

	pexDiscovery, err := pex.NewPEXDiscovery(h, bootstrapMaddrs, pexDiscoveryUpdateTime)
	if err != nil {
		return nil, xerrors.Errorf("failed to setup pex pexDiscovery: %v", err)
	}

	return pexDiscovery, nil
}

func Start() {
	configPath := flag.String("config", "", "Path to config")
	verbose := flag.Bool("verbose", false, "Verbose logging")
	flag.Parse()

	if *configPath == "" {
		logrus.Fatal("no config path provided")
	}
	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		logrus.Fatalf("failed to load config: %v", err)
	}

	var privateKey crypto.PrivKey

	if cfg.IsBootstrap {
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
		if err != nil {
			logrus.Fatal(err)
		}
	}

	node, err := NewNode(cfg, privateKey, DefaultPEXUpdateTime)
	if err != nil {
		logrus.Fatal(err)
	}

	// log
	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}

	//log.SetDebugLogging()

	//ctx, ctxCancel := context.WithCancel(context.Background())
	//node.GlobalCtx = ctx
	//node.GlobalCtxCancel = ctxCancel

	err = node.Run(context.TODO())
	if err != nil {
		logrus.Fatal(err)
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

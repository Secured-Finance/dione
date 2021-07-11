package node

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/multiformats/go-multiaddr"

	"github.com/asaskevich/EventBus"

	"github.com/fxamacker/cbor/v2"

	"github.com/Secured-Finance/dione/types"

	"github.com/Secured-Finance/dione/blockchain"

	types2 "github.com/Secured-Finance/dione/blockchain/types"

	gorpc "github.com/libp2p/go-libp2p-gorpc"

	"github.com/Secured-Finance/dione/blockchain/pool"

	"github.com/Secured-Finance/dione/blockchain/sync"

	"github.com/Secured-Finance/dione/consensus"
	pubsub2 "github.com/Secured-Finance/dione/pubsub"

	"github.com/libp2p/go-libp2p-core/discovery"

	"github.com/Secured-Finance/dione/rpc"
	rtypes "github.com/Secured-Finance/dione/rpc/types"

	solana2 "github.com/Secured-Finance/dione/rpc/solana"

	"github.com/Secured-Finance/dione/rpc/filecoin"

	"golang.org/x/xerrors"

	"github.com/Secured-Finance/dione/beacon"

	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/ethclient"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
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
	DisputeManager   *consensus.DisputeManager
	BlockPool        *blockchain.BlockChain
	MemPool          *pool.Mempool
	SyncManager      sync.SyncManager
	NetworkService   *NetworkService
	NetworkRPCHost   *gorpc.Server
	Bus              EventBus.Bus
	//Cache            cache.Cache
	//Wallet           *wallet.LocalWallet
}

func NewNode(config *config.Config, prvKey crypto.PrivKey, pexDiscoveryUpdateTime time.Duration) (*Node, error) {
	n := &Node{
		Config: config,
	}

	bus := EventBus.New()
	n.Bus = bus

	// initialize libp2p host
	lhost, err := provideLibp2pHost(n.Config, prvKey)
	if err != nil {
		logrus.Fatal(err)
	}
	n.Host = lhost
	logrus.Info("Libp2p host has been successfully initialized!")

	// initialize ethereum client
	ethClient, err := provideEthereumClient(n.Config)
	if err != nil {
		logrus.Fatal(err)
	}
	n.Ethereum = ethClient
	logrus.Info("Ethereum client has been successfully initialized!")

	// initialize blockchain rpc clients
	err = n.setupRPCClients()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("RPC clients has been successfully configured!")

	// initialize pubsub subsystem
	psb := providePubsubRouter(lhost, n.Config)
	n.PubSubRouter = psb
	logrus.Info("PubSub subsystem has been initialized!")

	// get list of bootstrap multiaddresses
	baddrs, err := provideBootstrapAddrs(n.Config)
	if err != nil {
		logrus.Fatal(err)
	}

	// initialize peer discovery
	peerDiscovery, err := providePeerDiscovery(baddrs, lhost, pexDiscoveryUpdateTime)
	if err != nil {
		logrus.Fatal(err)
	}
	n.PeerDiscovery = peerDiscovery
	logrus.Info("Peer discovery subsystem has been initialized!")

	// initialize event log cache subsystem
	//c := provideCache(config)
	//n.Cache = c
	//logrus.Info("Event cache subsystem has initialized!")

	// == initialize blockchain modules

	// initialize blockpool database
	bc, err := provideBlockChain(n.Config)
	if err != nil {
		logrus.Fatalf("Failed to initialize blockpool: %s", err.Error())
	}
	n.BlockPool = bc
	logrus.Info("Block pool database has been successfully initialized!")

	// initialize mempool
	mp, err := provideMemPool()
	if err != nil {
		logrus.Fatalf("Failed to initialize mempool: %s", err.Error())
	}
	n.MemPool = mp
	logrus.Info("Mempool has been successfully initialized!")

	ns := provideNetworkService(bc)
	n.NetworkService = ns
	rpcHost := provideNetworkRPCHost(lhost)
	err = rpcHost.Register(ns)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Node p2p RPC network service has been successfully initialized!")

	// initialize libp2p-gorpc client
	r := provideP2PRPCClient(lhost)

	// initialize sync manager

	var baddr multiaddr.Multiaddr
	if len(baddrs) == 0 {
		baddr = nil
	} else {
		baddr = baddrs[0]
	}
	sm, err := provideSyncManager(bus, bc, mp, r, baddr, psb) // FIXME here we just pick up first bootstrap in list
	if err != nil {
		logrus.Fatal(err)
	}
	n.SyncManager = sm
	logrus.Info("Blockchain synchronization subsystem has been successfully initialized!")

	// initialize mining subsystem
	miner := provideMiner(n.Host.ID(), *n.Ethereum.GetEthAddress(), n.Ethereum, prvKey, mp)
	n.Miner = miner
	logrus.Info("Mining subsystem has been initialized!")

	// initialize consensus subsystem
	consensusManager := provideConsensusManager(bus, psb, miner, bc, ethClient, prvKey, n.Config.ConsensusMinApprovals)
	n.ConsensusManager = consensusManager
	logrus.Info("Consensus subsystem has been initialized!")

	// initialize random beacon network subsystem
	randomBeaconNetwork, err := provideBeacon(psb.Pubsub, consensusManager)
	if err != nil {
		logrus.Fatal(err)
	}
	n.Beacon = randomBeaconNetwork
	logrus.Info("Random beacon subsystem has been initialized!")

	// initialize dispute subsystem
	disputeManager, err := provideDisputeManager(context.TODO(), ethClient, consensusManager, config, bc)
	if err != nil {
		logrus.Fatal(err)
	}
	n.DisputeManager = disputeManager
	logrus.Info("Dispute subsystem has been initialized!")

	// initialize internal eth wallet
	//w, err := provideWallet(n.Host.ID(), rawPrivKey)
	//if err != nil {
	//	logrus.Fatal(err)
	//}
	//n.Wallet = w

	return n, nil
}

func (n *Node) Run(ctx context.Context) error {
	err := n.runLibp2pAsync(ctx)
	if err != nil {
		return err
	}
	n.subscribeOnEthContractsAsync(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
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
					rpcMethod := rpc.GetRPCMethod(event.OriginChain, event.RequestType)
					if rpcMethod == nil {
						logrus.Errorf("Invalid RPC method name/type %d/%s for oracle request %s", event.OriginChain, event.RequestType, event.ReqID.String())
						continue
					}
					res, err := rpcMethod(event.RequestParams)
					if err != nil {
						logrus.Errorf("Failed to invoke RPC method for oracle request %s: %s", event.ReqID.String(), err.Error())
						continue
					}
					task := &types.DioneTask{
						OriginChain:   event.OriginChain,
						RequestType:   event.RequestType,
						RequestParams: event.RequestParams,
						Payload:       res,
						RequestID:     event.ReqID.String(),
					}
					data, err := cbor.Marshal(task)
					if err != nil {
						logrus.Errorf("Failed to marshal RPC response for oracle request %s: %s", event.ReqID.String(), err.Error())
						continue
					}
					tx := types2.CreateTransaction(data)
					err = n.MemPool.StoreTx(tx)
					if err != nil {
						logrus.Errorf("Failed to store tx in mempool: %s", err.Error())
						continue
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

func (n *Node) setupRPCClients() error {
	fc := filecoin.NewLotusClient()
	rpc.RegisterRPC(rtypes.RPCTypeFilecoin, map[string]func(string) ([]byte, error){
		"getTransaction": fc.GetTransaction,
		"getBlock":       fc.GetBlock,
	})

	sl := solana2.NewSolanaClient()
	rpc.RegisterRPC(rtypes.RPCTypeSolana, map[string]func(string) ([]byte, error){
		"getTransaction": sl.GetTransaction,
	})

	return nil
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
		// FIXME just a little hack
		if _, err := os.Stat(".bootstrap_privkey"); os.IsNotExist(err) {
			privateKey, err = generatePrivateKey()
			if err != nil {
				logrus.Fatal(err)
			}

			f, _ := os.Create(".bootstrap_privkey")
			r, _ := privateKey.Raw()
			_, err = f.Write(r)
			if err != nil {
				logrus.Fatal(err)
			}
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

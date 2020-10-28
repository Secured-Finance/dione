package node

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht/dual"

	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/consensus"
	"github.com/Secured-Finance/dione/pb"
	"github.com/Secured-Finance/dione/rpc"
	"github.com/Secured-Finance/dione/rpcclient"
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
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

func (n *Node) setupNode(ctx context.Context, prvKey crypto.PrivKey) {
	n.setupLibp2pHost(context.TODO(), prvKey)
	//n.setupEthereumClient()
	//n.setupFilecoinClient()
	n.setupPubsub()
	n.setupConsensusManager()
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
	time.Sleep(3 * time.Second)
}

func (n *Node) setupConsensusManager() {
	n.ConsensusManager = consensus.NewPBFTConsensusManager(n.PubSubRouter, 2)
}

func (n *Node) setupLibp2pHost(ctx context.Context, privateKey crypto.PrivKey) {
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

	kademliaDHT, err := dual.New(context.Background(), n.Host)
	if err != nil {
		logrus.Fatal("Failed to create new DHT instance: ", err)
	}

	if err = kademliaDHT.Bootstrap(context.Background()); err != nil {
		logrus.Fatal(err)
	}

	if !n.Config.Bootstrap {
		var wg sync.WaitGroup
		bootstrapMultiaddr, err := multiaddr.NewMultiaddr(n.Config.BootstrapNodeMultiaddr)
		if err != nil {
			logrus.Fatal(err)
		}
		peerinfo, _ := peer.AddrInfoFromP2pAddr(bootstrapMultiaddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := n.Host.Connect(context.Background(), *peerinfo); err != nil {
				logrus.Fatal(err)
			}
			logrus.Info("Connection established with bootstrap node:", *peerinfo)
		}()
		wg.Wait()
	}

	logrus.Info("Announcing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(context.Background(), routingDiscovery, n.Config.Rendezvous)
	logrus.Info("Successfully announced!")

	// Randezvous string = service tag
	// Discover all peers with our service
	logrus.Info("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(context.Background(), n.Config.Rendezvous)
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
					logrus.Info("Found peer:", newPeer, ", put it to the peerstore")
					n.Host.Peerstore().AddAddr(newPeer.ID, newPeer.Addrs[0], peerstore.PermanentAddrTTL)
					// Connect to the peer
					if err := n.Host.Connect(ctx, newPeer); err != nil {
						logrus.Warn("Connection failed: ", err)
					}
					logrus.Info("Connected to: ", newPeer)
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

	node.setupNode(ctx, privKey)
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

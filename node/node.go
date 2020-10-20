package node

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"sync"

	"github.com/Secured-Finance/p2p-oracle-node/config"
	"github.com/Secured-Finance/p2p-oracle-node/rpc"
	"github.com/Secured-Finance/p2p-oracle-node/rpcclient"
	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
)

type Node struct {
	Host            host.Host
	PubSub          *pubsub.PubSub
	GlobalCtx       context.Context
	GlobalCtxCancel context.CancelFunc
	OracleTopic     string
	Config          *config.Config
	Logger          *log.ZapEventLogger
	Lotus           *rpc.LotusClient
	Ethereum        *rpcclient.EthereumClient
}

func NewNode(configPath string) (*Node, error) {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return nil, err
	}
	node := &Node{
		OracleTopic: "p2p_oracle",
		Config:      cfg,
		Logger:      log.Logger("node"),
	}
	log.SetAllLoggers(log.LevelInfo)

	return node, nil
}

func (n *Node) setupNode(ctx context.Context, prvKey crypto.PrivKey) {
	listenMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", n.Config.ListenAddr, n.Config.ListenPort))
	if err != nil {
		n.Logger.Fatal("Failed to generate new node multiaddress:", err)
	}
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(listenMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		n.Logger.Fatal("Failed to set a new libp2p node:", err)
	}
	n.Host = host
	n.bootstrapLibp2pHost(context.TODO())
	n.setupEthereumClient()
	n.setupFilecoinClient()
	//n.startPubSub(ctx, host)
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

func (n *Node) bootstrapLibp2pHost(ctx context.Context) {
	kademliaDHT, err := dht.New(context.Background(), n.Host)
	if err != nil {
		n.Logger.Fatal("Failed to create new DHT instance: ", err)
	}

	if err = kademliaDHT.Bootstrap(context.Background()); err != nil {
		n.Logger.Fatal(err)
	}

	if !n.Config.Bootstrap {
		var wg sync.WaitGroup
		bootstrapMultiaddr, err := multiaddr.NewMultiaddr(n.Config.BootstrapNodeMultiaddr)
		if err != nil {
			n.Logger.Fatal(err)
		}
		peerinfo, _ := peer.AddrInfoFromP2pAddr(bootstrapMultiaddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := n.Host.Connect(context.Background(), *peerinfo); err != nil {
				n.Logger.Fatal(err)
			}
			n.Logger.Info("Connection established with bootstrap node:", *peerinfo)
		}()
		wg.Wait()
	}

	n.Logger.Info("Announcing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(context.Background(), routingDiscovery, n.Config.Rendezvous)
	n.Logger.Info("Successfully announced!")

	// Randezvous string = service tag
	// Disvover all peers with our service (all ms devices)
	n.Logger.Info("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(context.Background(), n.Config.Rendezvous)
	if err != nil {
		n.Logger.Fatal("Failed to find new peers, exiting...", err)
	}
	go func() {
	MainLoop:
		for {
			select {
			case <-ctx.Done():
				break MainLoop
			case newPeer := <-peerChan:
				{
					n.Logger.Info("Found peer:", newPeer, ", put it to the peerstore")
					n.Host.Peerstore().AddAddr(newPeer.ID, newPeer.Addrs[0], peerstore.PermanentAddrTTL)
					// Connect to the peer
					if err := n.Host.Connect(ctx, newPeer); err != nil {
						n.Logger.Warn("Connection failed: ", err)
					}
					n.Logger.Info("Connected to: ", newPeer)
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
		log.SetAllLoggers(log.LevelDebug)
	} else {
		log.SetAllLoggers(log.LevelInfo)
	}
	if err != nil {
		log.Logger("node").Panic(err)
	}

	privKey, err := generatePrivateKey()
	if err != nil {
		node.Logger.Fatal(err)
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	node.GlobalCtx = ctx
	node.GlobalCtxCancel = ctxCancel

	node.setupNode(ctx, privKey)
	return nil
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

package node

import (
	"context"
	"fmt"
	"time"

	drand2 "github.com/Secured-Finance/dione/beacon/drand"

	"github.com/libp2p/go-libp2p-core/protocol"

	gorpc "github.com/libp2p/go-libp2p-gorpc"

	"github.com/Secured-Finance/dione/blockchain/sync"

	"github.com/Secured-Finance/dione/blockchain/pool"

	"github.com/Secured-Finance/dione/beacon"
	"github.com/Secured-Finance/dione/cache"
	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/consensus"
	"github.com/Secured-Finance/dione/ethclient"
	"github.com/Secured-Finance/dione/pubsub"
	"github.com/Secured-Finance/dione/types"
	"github.com/Secured-Finance/dione/wallet"
	pex "github.com/Secured-Finance/go-libp2p-pex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/discovery"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub2 "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/xerrors"
)

const (
	DioneProtocolID = protocol.ID("/dione/1.0")
)

func provideCache(config *config.Config) cache.Cache {
	var backend cache.Cache
	switch config.CacheType {
	case "in-memory":
		backend = cache.NewInMemoryCache()
	case "redis":
		backend = cache.NewRedisCache(config)
	default:
		backend = cache.NewInMemoryCache()
	}
	return backend
}

func provideDisputeManager(ctx context.Context, ethClient *ethclient.EthereumClient, pcm *consensus.PBFTConsensusManager, cfg *config.Config) (*consensus.DisputeManager, error) {
	return consensus.NewDisputeManager(ctx, ethClient, pcm, cfg.Ethereum.DisputeVoteWindow)
}

func provideMiner(peerID peer.ID, ethAddress common.Address, beacon beacon.BeaconNetworks, ethClient *ethclient.EthereumClient, privateKey []byte) *consensus.Miner {
	return consensus.NewMiner(peerID, ethAddress, beacon, ethClient, privateKey)
}

func provideBeacon(ps *pubsub2.PubSub) (beacon.BeaconNetworks, error) {
	networks := beacon.BeaconNetworks{}
	bc, err := drand2.NewDrandBeacon(ps)
	if err != nil {
		return nil, fmt.Errorf("failed to setup drand beacon: %w", err)
	}
	networks = append(networks, beacon.BeaconNetwork{Start: types.DrandRound(config.ChainGenesis), Beacon: bc})
	// NOTE: currently we use only one network
	return networks, nil
}

// FIXME: do we really need this?
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
	err := ethereum.Initialize(&config.Ethereum)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize ethereum client: %v", err)
	}
	return ethereum, nil
}

func providePubsubRouter(lhost host.Host, config *config.Config) *pubsub.PubSubRouter {
	return pubsub.NewPubSubRouter(lhost, config.PubSub.ServiceTopicName, config.IsBootstrap)
}

func provideConsensusManager(psb *pubsub.PubSubRouter, miner *consensus.Miner, ethClient *ethclient.EthereumClient, privateKey []byte, minApprovals int, evc cache.Cache) *consensus.PBFTConsensusManager {
	return consensus.NewPBFTConsensusManager(psb, minApprovals, privateKey, ethClient, miner, evc)
}

func provideLibp2pHost(config *config.Config, privateKey crypto.PrivKey) (host.Host, error) {
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

func provideNetworkRPCHost(h host.Host) *gorpc.Server {
	return gorpc.NewServer(h, DioneProtocolID)
}

func provideBootstrapAddrs(c *config.Config) ([]multiaddr.Multiaddr, error) {
	if c.IsBootstrap {
		return nil, nil
	}

	var bootstrapMaddrs []multiaddr.Multiaddr
	for _, a := range c.BootstrapNodes {
		maddr, err := multiaddr.NewMultiaddr(a)
		if err != nil {
			return nil, xerrors.Errorf("invalid multiaddress of bootstrap node: %v", err)
		}
		bootstrapMaddrs = append(bootstrapMaddrs, maddr)
	}

	return bootstrapMaddrs, nil
}

func providePeerDiscovery(baddrs []multiaddr.Multiaddr, h host.Host, pexDiscoveryUpdateTime time.Duration) (discovery.Discovery, error) {
	pexDiscovery, err := pex.NewPEXDiscovery(h, baddrs, pexDiscoveryUpdateTime)
	if err != nil {
		return nil, xerrors.Errorf("failed to setup pex pexDiscovery: %v", err)
	}

	return pexDiscovery, nil
}

func provideBlockPool(config *config.Config) (*pool.BlockPool, error) {
	return pool.NewBlockPool(config.Blockchain.DatabasePath)
}

func provideMemPool() (*pool.Mempool, error) {
	return pool.NewMempool()
}

func provideSyncManager(bp *pool.BlockPool, mp *pool.Mempool, r *gorpc.Client, bootstrap multiaddr.Multiaddr, psb *pubsub.PubSubRouter) (sync.SyncManager, error) {
	addr, err := peer.AddrInfoFromP2pAddr(bootstrap)
	if err != nil {
		return nil, err
	}
	return sync.NewSyncManager(bp, mp, r, addr.ID, psb), nil
}

func provideP2PRPCClient(h host.Host) *gorpc.Client {
	return gorpc.NewClient(h, DioneProtocolID)
}

func provideNetworkService(bp *pool.BlockPool) *NetworkService {
	return NewNetworkService(bp)
}

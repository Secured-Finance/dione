package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/google/logger"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	crypto "github.com/libp2p/go-libp2p-crypto"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	listenPort := flag.Int("port", 0, "Listen port number")
	listenAddr := flag.String("addr", "", "Listen address")
	verbose := flag.Bool("verbose", false, "Verbose logs")
	syslog := flag.Bool("syslog", false, "Log to system logging daemon")
	bootstrap := flag.Bool("bootstrap", false, "Start up bootstrap node")
	bootstrapAddress := flag.String("baddr", "", "Address of bootstrap node")
	rendezvousString := flag.String("rendezvous", "", "DHT rendezvous string")

	flag.Parse()

	defer logger.Init("node", *verbose, *syslog, ioutil.Discard).Close()

	r := rand.Reader

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		logger.Fatal(err)
	}

	ctx := context.Background()

	listenMaddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", *listenAddr, *listenPort))
	if err != nil {
		logger.Fatal(err)
	}
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(listenMaddr),
		libp2p.Identity(prvKey),
	)

	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		panic(err)
	}

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		logger.Fatal(err)
	}

	if !*bootstrap {
		var wg sync.WaitGroup
		bMaddr, err := multiaddr.NewMultiaddr(*bootstrapAddress)
		if err != nil {
			logger.Fatal(err)
		}
		peerinfo, _ := peer.AddrInfoFromP2pAddr(bMaddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				logger.Fatal(err)
			} else {
				logger.Info("Connection established with bootstrap node:", *peerinfo)
			}
		}()
		wg.Wait()
	}

	logger.Info("Libp2p node is successfully started!")
	multiaddress := fmt.Sprintf("/ip4/%s/tcp/%v/p2p/%s", *listenAddr, *listenPort, host.ID().Pretty())
	logger.Infof("Your multiaddress: %s", multiaddress)

	logger.Info("Announcing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(ctx, routingDiscovery, *rendezvousString)
	logger.Info("Successfully announced!")

	logger.Info("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(ctx, *rendezvousString)
	if err != nil {
		logger.Fatal(err)
	}

	for peer := range peerChan {
		if peer.ID == host.ID() {
			continue
		}
		logger.Info("Found peer:", peer)

		logger.Info("Connecting to:", peer)
		err := host.Connect(ctx, peer)
		if err != nil {
			logger.Error("Connection failed: " + err.Error())
			continue
		}

		logger.Info("Connected to: ", peer)
	}

	select {}
}

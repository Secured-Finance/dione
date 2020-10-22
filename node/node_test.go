package node

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Secured-Finance/dione/config"
	"github.com/sirupsen/logrus"
)

func TestConsensus(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	//log.SetAllLoggers(log.LevelDebug)

	cfg := &config.Config{
		ListenPort: "1234",
		ListenAddr: "0.0.0.0",
		Bootstrap:  true,
		Rendezvous: "dione",
		PubSub: config.PubSubConfig{
			ProtocolID: "/test/1.0",
		},
	}

	//cfg.BootstrapNodeMultiaddr = "/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN"

	// setup first node
	node1 := newNode(cfg)

	// setup second node
	cfg.ListenPort = "1235"
	cfg.Bootstrap = false
	cfg.BootstrapNodeMultiaddr = node1.Host.Addrs()[0].String() + fmt.Sprintf("/p2p/%s", node1.Host.ID().String())
	node2 := newNode(cfg)

	// setup third node
	cfg.ListenPort = "1236"
	node3 := newNode(cfg)

	cfg.ListenPort = "1237"
	node4 := newNode(cfg)
	cfg.ListenPort = "1238"
	node5 := newNode(cfg)
	cfg.ListenPort = "1239"
	node6 := newNode(cfg)


	time.Sleep(10 * time.Second)
	go node2.ConsensusManager.NewTestConsensus("test", "123")
	go node1.ConsensusManager.NewTestConsensus("test1", "123")
	go node3.ConsensusManager.NewTestConsensus("test", "123")
	go node4.ConsensusManager.NewTestConsensus("test1", "123")
	go node5.ConsensusManager.NewTestConsensus("test", "123")
	go node6.ConsensusManager.NewTestConsensus("test2", "123")
	select{}
}

func newNode(cfg *config.Config) *Node {
	privKey, err := generatePrivateKey()
	if err != nil {
		logrus.Fatal(err)
	}
	ctx, ctxCancel := context.WithCancel(context.Background())

	node := &Node{
		OracleTopic:     "dione",
		Config:          cfg,
		GlobalCtx:       ctx,
		GlobalCtxCancel: ctxCancel,
	}
	node.setupNode(ctx, privKey)
	return node
}
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
	privKey, err := generatePrivateKey()
	if err != nil {
		t.Error(err)
	}
	ctx, ctxCancel := context.WithCancel(context.Background())
	node1 := &Node{
		OracleTopic:     "dione",
		Config:          cfg,
		GlobalCtx:       ctx,
		GlobalCtxCancel: ctxCancel,
	}
	node1.setupNode(ctx, privKey)

	// setup second node
	privKey, err = generatePrivateKey()
	if err != nil {
		t.Error(err)
	}
	ctx, ctxCancel = context.WithCancel(context.Background())
	cfg.ListenPort = "1235"
	cfg.Bootstrap = false
	cfg.BootstrapNodeMultiaddr = node1.Host.Addrs()[0].String() + fmt.Sprintf("/p2p/%s", node1.Host.ID().String())

	node2 := &Node{
		OracleTopic:     "dione",
		Config:          cfg,
		GlobalCtx:       ctx,
		GlobalCtxCancel: ctxCancel,
	}
	node2.setupNode(ctx, privKey)

	// setup third node
	privKey, err = generatePrivateKey()
	if err != nil {
		t.Error(err)
	}
	ctx, ctxCancel = context.WithCancel(context.Background())
	cfg.ListenPort = "1236"
	node3 := &Node{
		OracleTopic:     "dione",
		Config:          cfg,
		GlobalCtx:       ctx,
		GlobalCtxCancel: ctxCancel,
	}
	node3.setupNode(ctx, privKey)


	time.Sleep(10 * time.Second)
	go node2.ConsensusManager.NewTestConsensus("test", "123")
	go node1.ConsensusManager.NewTestConsensus("test1", "123")
	go node3.ConsensusManager.NewTestConsensus("test", "123")
	select{}
}
package node

import (
	"context"
	"fmt"
	"testing"

	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/consensus"
	"github.com/ipfs/go-log"
)

func TestConsensus(t *testing.T) {
	var logger = log.Logger("test")
	log.SetAllLoggers(log.LevelDebug)

	cfg := &config.Config{
		ListenPort: "1234",
		ListenAddr: "127.0.0.1",
		Bootstrap:  true,
		Rendezvous: "dione",
		PubSub: config.PubSubConfig{
			ProtocolID: "/test/1.0",
		},
	}

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
		Logger:          log.Logger("node"),
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
		Logger:          log.Logger("node"),
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
		Logger:          log.Logger("node"),
	}
	node3.setupNode(ctx, privKey)

	logger.Debug(node3.Host.Peerstore().Peers())

	node2.ConsensusManager.NewTestConsensus("test")
	node1.ConsensusManager.NewTestConsensus("test1")
	node3.ConsensusManager.NewTestConsensus("test")
	var last consensus.ConsensusState = -1
	for {
		for _, v := range node1.ConsensusManager.Consensuses {
			if v.State != last {
				last = v.State
				t.Log("new state: " + fmt.Sprint(v.State))
			}
		}
	}
}

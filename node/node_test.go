package node

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Secured-Finance/dione/config"
	"github.com/sirupsen/logrus"
)

func TestConsensus(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	//log.SetAllLoggers(log.LevelDebug)

	// boolgen := newBoolgen()
	rand.Seed(time.Now().UnixNano())
	port := rand.Intn(100) + 10000

	cfg := &config.Config{
		ListenPort: port,
		ListenAddr: "0.0.0.0",
		Rendezvous: "dione",
		PubSub: config.PubSubConfig{
			ProtocolID: "/dione/1.0",
		},
		ConsensusMinApprovals: 3,
	}

	var nodes []*Node

	bNode := newNode(cfg)
	t.Logf("Bootstrap ID: %s", bNode.Host.ID())
	cfg.BootstrapNodes = []string{bNode.Host.Addrs()[0].String() + fmt.Sprintf("/p2p/%s", bNode.Host.ID().String())}
	nodes = append(nodes, bNode)

	maxNodes := 10

	for i := 1; i <= maxNodes; i++ {
		cfg.ListenPort += 1
		node := newNode(cfg)
		nodes = append(nodes, node)
	}

	time.Sleep(5 * time.Second)

	// var wg sync.WaitGroup

	// wg.Add(len(nodes))
	// for _, n := range nodes {
	// 	var testData string
	// 	if boolgen.Bool() {
	// 		testData = "test"
	// 	} else {
	// 		testData = "test1"
	// 	}
	// 	n.ConsensusManager.NewTestConsensus(testData, "123", func(finalData string) {
	// 		if finalData != "test" {
	// 			t.Errorf("Expected final data %s, not %s", "test", finalData)
	// 		}
	// 		wg.Done()
	// 	})
	// }
	// wg.Wait()
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
	node.setupNode(ctx, privKey, 1*time.Second)
	return node
}

type boolgen struct {
	src       rand.Source
	cache     int64
	remaining int
}

func newBoolgen() *boolgen {
	return &boolgen{src: rand.NewSource(time.Now().UnixNano())}
}

func (b *boolgen) Bool() bool {
	if b.remaining == 0 {
		b.cache, b.remaining = b.src.Int63(), 63
	}

	result := b.cache&0x01 == 1
	b.cache >>= 1
	b.remaining--

	return result
}

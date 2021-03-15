package consensus

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/node"
	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/sirupsen/logrus"
)

func TestConsensus(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	// setting up nodes
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

	var nodes []*node.Node

	bNode := newNode(cfg)
	t.Logf("Bootstrap ID: %s", bNode.Host.ID())
	cfg.BootstrapNodes = []string{bNode.Host.Addrs()[0].String() + fmt.Sprintf("/p2p/%s", bNode.Host.ID().String())}
	nodes = append(nodes, bNode)

	maxNodes := 10

	for i := 1; i <= maxNodes; i++ {
		cfg.ListenPort++
		node := newNode(cfg)
		nodes = append(nodes, node)
	}

	time.Sleep(5 * time.Second)

}

func newNode(cfg *config.Config) *node.Node {
	privKey, err := generatePrivateKey()
	if err != nil {
		logrus.Fatal(err)
	}

	node, err := node.NewNode(cfg, privKey, 1*time.Second)
	if err != nil {
		logrus.Fatal(err)
	}
	return node
}

func generatePrivateKey() (crypto.PrivKey, error) {
	r := crand.Reader
	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, r)
	if err != nil {
		return nil, err
	}
	return prvKey, nil
}

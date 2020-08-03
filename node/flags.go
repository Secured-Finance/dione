package node

import (
	"flag"

	"github.com/Secured-Finance/p2p-oracle-node/config"
)

func (node *Node) parseFlags() {
	listenPort := flag.String("port", node.Config.ListenPort, "Listen port number")
	listenAddr := flag.String("addr", node.Config.ListenAddr, "Listen address")
	bootstrap := flag.Bool("bootstrap", node.Config.Bootstrap, "Start up bootstrap node")
	bootstrapAddress := flag.String("baddr", node.Config.BootstrapNodeMultiaddr, "Address of bootstrap node")
	rendezvousString := flag.String("rendezvous", node.Config.Rendezvous, "DHT rendezvous string")
	protocolID := flag.String("protocol-id", node.Config.ProtocolID, "PubSub protocol ID")

	flag.Parse()

	new_config := &config.Config{
		ListenPort:             *listenPort,
		ListenAddr:             *listenAddr,
		Bootstrap:              *bootstrap,
		BootstrapNodeMultiaddr: *bootstrapAddress,
		Rendezvous:             *rendezvousString,
		ProtocolID:             *protocolID,
	}
	node.Config = new_config
}

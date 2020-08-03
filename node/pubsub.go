package node

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/Secured-Finance/p2p-oracle-node/handler"
	"github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
)

func (node *Node) readSub(subscription *pubsub.Subscription, incomingMessagesChan chan pubsub.Message) {
	ctx := node.GlobalCtx
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		msg, err := subscription.Next(context.Background())
		if err != nil {
			node.Logger.Warn("Error reading from buffer", err)
			return
		}

		if string(msg.Data) == "" {
			return
		}
		if string(msg.Data) != "\n" {
			addr, err := peer.IDFromBytes(msg.From)
			if err != nil {
				node.Logger.Warn("Error occurred when reading message From field...", err)
				return
			}

			// This checks if sender address of incoming message is ours. It is need because we get our messages when subscribed to the same topic.
			if addr == node.Host.ID() {
				continue
			}
			incomingMessagesChan <- *msg
		}

	}
}

// Subscribes to a topic and then get messages ..
func (node *Node) newTopic(topic string) {
	ctx := node.GlobalCtx
	subscription, err := node.PubSub.Subscribe(topic)
	if err != nil {
		node.Logger.Warn("Error occurred when subscribing to topic", err)
		return
	}
	time.Sleep(3 * time.Second)
	incomingMessages := make(chan pubsub.Message)

	go node.readSub(subscription, incomingMessages)
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-incomingMessages:
			{
				node.handler.HandleIncomingMessage(node.OracleTopic, msg, func(textMessage handler.TextMessage) {
					node.Logger.Info("%s \x1b[32m%s\x1b[0m> ", textMessage.From, textMessage.Body)
				})
			}
		}
	}
}

// Write messages to subscription (topic)
// NOTE: we don't need to be subscribed to publish something
func (node *Node) writeTopic(topic string) {
	ctx := node.GlobalCtx
	stdReader := bufio.NewReader(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		node.Logger.Info("> ")
		text, err := stdReader.ReadString('\n')
		if err != nil {

			if err == io.EOF {
				break
			}

			node.Logger.Warn("Error reading from stdin", err)
			return
		}
		message := &handler.BaseMessage{
			Body: text,
			Flag: handler.FlagGenericMessage,
		}

		sendData, err := json.Marshal(message)
		if err != nil {
			node.Logger.Warn("Error occurred when marshalling message object")
			continue
		}
		err = node.PubSub.Publish(topic, sendData)
		if err != nil {
			node.Logger.Warn("Error occurred when publishing", err)
			return
		}
	}
}

func (node *Node) getNetworkTopics() {
	// ctx := node.GlobalCtx
	node.handler.RequestNetworkTopics()
}

func (node *Node) startPubSub(ctx context.Context, host host.Host) {
	pb, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		node.Logger.Fatal("Error occurred when create PubSub", err)
	}

	// pb, err := pubsub.NewFloodsubWithProtocols(context.Background(), host, []protocol.ID{protocol.ID(node.Config.ProtocolID)}, pubsub.WithMessageSigning(true), pubsub.WithStrictSignatureVerification(true))
	// if err != nil {
	// 	node.Logger.Fatal("Error occurred when create PubSub", err)
	// }

	// Set global PubSub object
	node.PubSub = pb

	node.handler = handler.NewHandler(pb, node.OracleTopic, host.ID(), node.networkTopics)

	kademliaDHT, err := dht.New(ctx, host)
	if err != nil {
		node.Logger.Fatal("Failed to set a new DHT:", err)
	}

	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		node.Logger.Fatal(err)
	}

	if !node.Config.Bootstrap {
		var wg sync.WaitGroup
		bootstrapMultiaddr, err := multiaddr.NewMultiaddr(node.Config.BootstrapNodeMultiaddr)
		if err != nil {
			node.Logger.Fatal(err)
		}
		peerinfo, _ := peer.AddrInfoFromP2pAddr(bootstrapMultiaddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				node.Logger.Fatal(err)
			} else {
				node.Logger.Info("Connection established with bootstrap node:", *peerinfo)
			}
		}()
		wg.Wait()
	}

	node.Logger.Info("Announcing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	discovery.Advertise(ctx, routingDiscovery, node.Config.Rendezvous)
	node.Logger.Info("Successfully announced!")

	// Randezvous string = service tag
	// Disvover all peers with our service (all ms devices)
	node.Logger.Info("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(ctx, node.Config.Rendezvous)
	if err != nil {
		node.Logger.Fatal("Failed to find new peers, exiting...", err)
	}

	// NOTE:  here we use Randezvous string as 'topic' by default .. topic != service tag
	node.OracleTopic = node.Config.Rendezvous
	subscription, err := pb.Subscribe(node.OracleTopic)
	if err != nil {
		node.Logger.Warn("Error occurred when subscribing to topic", err)
		return
	}

	node.Logger.Info("Waiting for correct set up of PubSub...")
	time.Sleep(3 * time.Second)

	incomingMessages := make(chan pubsub.Message)

	go func() {
		node.writeTopic(node.OracleTopic)
		node.GlobalCtxCancel()
	}()
	go node.readSub(subscription, incomingMessages)
	go node.getNetworkTopics()

MainLoop:
	for {
		select {
		case <-ctx.Done():
			break MainLoop
		case msg := <-incomingMessages:
			{
				node.handler.HandleIncomingMessage(node.OracleTopic, msg, func(textMessage handler.TextMessage) {
					node.Logger.Info("%s > \x1b[32m%s\x1b[0m", textMessage.From, textMessage.Body)
					node.Logger.Info("> ")
				})
			}
		case newPeer := <-peerChan:
			{
				node.Logger.Info("\nFound peer:", newPeer, ", add address to peerstore")

				// Adding peer addresses to local peerstore
				host.Peerstore().AddAddr(newPeer.ID, newPeer.Addrs[0], peerstore.PermanentAddrTTL)
				// Connect to the peer
				if err := host.Connect(ctx, newPeer); err != nil {
					node.Logger.Warn("Connection failed:", err)
				}
				node.Logger.Info("Connected to:", newPeer)
				node.Logger.Info("> ")
			}
		}
	}

	if err := host.Close(); err != nil {
		node.Logger.Info("\nClosing host failed:", err)
	}
	node.Logger.Info("\nBye")
}

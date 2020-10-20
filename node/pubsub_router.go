package node

import (
	"context"
	"encoding/json"

	"github.com/Secured-Finance/p2p-oracle-node/models"
	"github.com/ipfs/go-log"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type PubSubRouter struct {
	node          *Node
	pubsub        *pubsub.PubSub
	logger        *log.ZapEventLogger
	context       context.Context
	contextCancel context.CancelFunc
	handlers      map[string][]Handler
}

func NewPubSubRouter(n *Node) *PubSubRouter {
	ctx, ctxCancel := context.WithCancel(context.Background())

	psr := &PubSubRouter{
		node:          n,
		logger:        log.Logger("PubSubRouter"),
		context:       ctx,
		contextCancel: ctxCancel,
	}

	pb, err := pubsub.NewGossipSub(
		context.Background(),
		psr.node.Host, pubsub.WithMessageSigning(true),
		pubsub.WithStrictSignatureVerification(true),
	)
	if err != nil {
		psr.logger.Fatal("Error occurred when create PubSub", err)
	}

	n.OracleTopic = n.Config.Rendezvous
	subscription, err := pb.Subscribe(n.OracleTopic)
	if err != nil {
		psr.logger.Fatal("Error occurred when subscribing to service topic", err)
	}
	psr.pubsub = pb

	go func() {
		for {
			select {
			case <-psr.context.Done():
				return
			default:
				{
					msg, err := subscription.Next(psr.context)
					if err != nil {
						psr.logger.Warn("Failed to receive pubsub message: ", err.Error())
					}
					psr.handleMessage(msg)
				}
			}
		}
	}()

	return psr
}

func (psr *PubSubRouter) handleMessage(p *pubsub.Message) {
	senderPeerID, err := peer.IDFromBytes(p.From)
	if err != nil {
		psr.logger.Warn("Unable to decode sender peer ID! " + err.Error())
		return
	}
	// We can receive our own messages when sending to the topic. So we should drop them.
	if senderPeerID == psr.node.Host.ID() {
		return
	}
	var message models.Message
	err = json.Unmarshal(p.Data, &message)
	if err != nil {
		psr.logger.Warn("Unable to decode message data! " + err.Error())
		return
	}
	message.From = senderPeerID.String()
	handlers, ok := psr.handlers[message.Type]
	if !ok {
		psr.logger.Warn("Dropping message " + message.Type + " because we don't have any handlers!")
		return
	}
	for _, v := range handlers {
		go v.HandleMessage(&message)
	}
}

func (psr *PubSubRouter) Hook(messageType string, handler Handler) {
	handlers, ok := psr.handlers[messageType]
	if !ok {
		emptyArray := []Handler{}
		psr.handlers[messageType] = emptyArray
		handlers = emptyArray
	}
	psr.handlers[messageType] = append(handlers, handler)
}

func (psr *PubSubRouter) Shutdown() {
	psr.contextCancel()
}

package pb

import (
	"context"
	"encoding/json"

	"github.com/Secured-Finance/dione/models"
	host "github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

type PubSubRouter struct {
	node          host.Host
	Pubsub        *pubsub.PubSub
	context       context.Context
	contextCancel context.CancelFunc
	handlers      map[string][]Handler
	oracleTopic   string
}

func NewPubSubRouter(h host.Host, oracleTopic string) *PubSubRouter {
	ctx, ctxCancel := context.WithCancel(context.Background())

	psr := &PubSubRouter{
		node:          h,
		context:       ctx,
		contextCancel: ctxCancel,
		handlers:      make(map[string][]Handler),
	}

	pb, err := pubsub.NewGossipSub(
		context.TODO(),
		psr.node, //pubsub.WithMessageSigning(true),
		//pubsub.WithStrictSignatureVerification(true),
	)
	if err != nil {
		logrus.Fatal("Error occurred when create PubSub", err)
	}

	psr.oracleTopic = oracleTopic
	subscription, err := pb.Subscribe(oracleTopic)
	if err != nil {
		logrus.Fatal("Error occurred when subscribing to service topic", err)
	}
	psr.Pubsub = pb

	go func() {
		for {
			select {
			case <-psr.context.Done():
				return
			default:
				{
					msg, err := subscription.Next(psr.context)
					if err != nil {
						logrus.Warn("Failed to receive pubsub message: ", err.Error())
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
		logrus.Warn("Unable to decode sender peer ID! " + err.Error())
		return
	}
	// We can receive our own messages when sending to the topic. So we should drop them.
	if senderPeerID == psr.node.ID() {
		logrus.Debug("Drop message because it came from the current node - a bug (or feature) in the pubsub system")
		return
	}
	var message models.Message
	err = json.Unmarshal(p.Data, &message)
	if err != nil {
		logrus.Warn("Unable to decode message data! " + err.Error())
		return
	}
	message.From = senderPeerID.String()
	handlers, ok := psr.handlers[message.Type]
	if !ok {
		logrus.Warn("Dropping message " + message.Type + " because we don't have any handlers!")
		return
	}
	for _, v := range handlers {
		go v(&message)
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

func (psr *PubSubRouter) BroadcastToServiceTopic(msg *models.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = psr.Pubsub.Publish(psr.oracleTopic, data)
	return err
}

func (psr *PubSubRouter) Shutdown() {
	psr.contextCancel()
}

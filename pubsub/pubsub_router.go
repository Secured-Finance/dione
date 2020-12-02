package pubsub

import (
	"context"

	"github.com/fxamacker/cbor/v2"

	"github.com/Secured-Finance/dione/consensus/types"

	host "github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

type PubSubRouter struct {
	node                host.Host
	Pubsub              *pubsub.PubSub
	context             context.Context
	contextCancel       context.CancelFunc
	serviceSubscription *pubsub.Subscription
	handlers            map[types.MessageType][]Handler
	oracleTopic         string
}

func NewPubSubRouter(h host.Host, oracleTopic string) *PubSubRouter {
	ctx, ctxCancel := context.WithCancel(context.Background())

	psr := &PubSubRouter{
		node:          h,
		context:       ctx,
		contextCancel: ctxCancel,
		handlers:      make(map[types.MessageType][]Handler),
	}

	pb, err := pubsub.NewFloodSub(
		context.TODO(),
		psr.node, //pubsub.WithMessageSigning(true),
		//pubsub.WithStrictSignatureVerification(true),
	)
	if err != nil {
		logrus.Fatal("Error occurred when create PubSub", err)
	}

	psr.oracleTopic = oracleTopic
	topic, err := pb.Join(oracleTopic)
	if err != nil {
		logrus.Fatal("Error occurred when subscribing to service topic", err)
	}

	subscription, err := topic.Subscribe()
	psr.serviceSubscription = subscription
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
		return
	}
	var message types.Message
	err = cbor.Unmarshal(p.Data, &message)
	if err != nil {
		logrus.Warn("Unable to decode message data! " + err.Error())
		return
	}
	message.From = senderPeerID
	handlers, ok := psr.handlers[message.Type]
	if !ok {
		logrus.Warn("Dropping message " + string(message.Type) + " because we don't have any handlers!")
		return
	}
	for _, v := range handlers {
		go v(&message)
	}
}

func (psr *PubSubRouter) Hook(messageType types.MessageType, handler Handler) {
	_, ok := psr.handlers[messageType]
	if !ok {
		psr.handlers[messageType] = []Handler{}
	}
	psr.handlers[messageType] = append(psr.handlers[messageType], handler)
}

func (psr *PubSubRouter) BroadcastToServiceTopic(msg *types.Message) error {
	data, err := cbor.Marshal(msg)
	if err != nil {
		return err
	}
	err = psr.Pubsub.Publish(psr.oracleTopic, data)
	return err
}

func (psr *PubSubRouter) Shutdown() {
	psr.contextCancel()
}

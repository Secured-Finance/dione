package pubsub

import (
	"context"
	"time"

	"github.com/fxamacker/cbor/v2"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sirupsen/logrus"
)

type PubSubRouter struct {
	node                host.Host
	Pubsub              *pubsub.PubSub
	context             context.Context
	contextCancel       context.CancelFunc
	serviceSubscription *pubsub.Subscription
	handlers            map[PubSubMessageType][]Handler
	oracleTopicName     string
	oracleTopic         *pubsub.Topic
	typeMapping         map[PubSubMessageType]interface{} // message type -> sample
}

type Handler func(message *GenericMessage)

func NewPubSubRouter(h host.Host, oracleTopic string, isBootstrap bool) *PubSubRouter {
	ctx, ctxCancel := context.WithCancel(context.Background())

	psr := &PubSubRouter{
		node:          h,
		context:       ctx,
		contextCancel: ctxCancel,
		handlers:      make(map[PubSubMessageType][]Handler),
		typeMapping:   map[PubSubMessageType]interface{}{},
	}

	var pbOptions []pubsub.Option

	if isBootstrap {
		// turn off the mesh in bootstrappers -- only do gossip and PX
		pubsub.GossipSubD = 0
		pubsub.GossipSubDscore = 0
		pubsub.GossipSubDlo = 0
		pubsub.GossipSubDhi = 0
		pubsub.GossipSubDout = 0
		pubsub.GossipSubDlazy = 64
		pubsub.GossipSubGossipFactor = 0.25
		pubsub.GossipSubPruneBackoff = 5 * time.Minute
		// turn on PX
		pbOptions = append(pbOptions, pubsub.WithPeerExchange(true))
	}

	pb, err := pubsub.NewGossipSub(
		context.TODO(),
		psr.node,
		pbOptions...,
	)

	if err != nil {
		logrus.Fatalf("Error occurred when initializing PubSub subsystem: %v", err)
	}

	psr.oracleTopicName = oracleTopic
	topic, err := pb.Join(oracleTopic)
	if err != nil {
		logrus.Fatalf("Error occurred when subscribing to service topic: %v", err)
	}

	subscription, err := topic.Subscribe()
	psr.serviceSubscription = subscription
	psr.Pubsub = pb
	psr.oracleTopic = topic

	go func() {
		for {
			select {
			case <-psr.context.Done():
				return
			default:
				{
					msg, err := subscription.Next(psr.context)
					if err != nil {
						logrus.Warnf("Failed to receive pubsub message: %v", err)
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
	var genericMessage PubSubMessage
	var message GenericMessage
	err = cbor.Unmarshal(p.Data, &genericMessage)
	if err != nil {
		logrus.Warn("Unable to decode pubsub message data! " + err.Error())
		return
	}
	sampleMsg, ok := psr.typeMapping[genericMessage.Type]
	if !ok {
		logrus.Warnf("Unknown message type %d: we have no clue how to decode it", genericMessage.Type)
		return
	}
	destMsg := sampleMsg
	err = cbor.Unmarshal(genericMessage.Payload, &destMsg)
	if err != nil {
		logrus.Warn("Unable to decode pubsub message data! " + err.Error())
		return
	}
	message.From = senderPeerID
	message.Type = genericMessage.Type
	message.Payload = destMsg
	handlers, ok := psr.handlers[genericMessage.Type]
	if !ok {
		logrus.Warn("Dropping pubsub message " + string(genericMessage.Type) + " because we don't have any handlers!")
		return
	}
	for _, v := range handlers {
		go v(&message)
	}
}

func (psr *PubSubRouter) Hook(messageType PubSubMessageType, handler Handler, sample interface{}) {
	_, ok := psr.handlers[messageType]
	if !ok {
		psr.handlers[messageType] = []Handler{}
	}
	psr.handlers[messageType] = append(psr.handlers[messageType], handler)
	psr.typeMapping[messageType] = sample
}

func (psr *PubSubRouter) BroadcastToServiceTopic(msg *GenericMessage) error {
	data, err := cbor.Marshal(msg)
	if err != nil {
		return err
	}
	err = psr.oracleTopic.Publish(context.TODO(), data)
	return err
}

func (psr *PubSubRouter) Shutdown() {
	psr.contextCancel()
}

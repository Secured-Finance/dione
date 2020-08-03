package handler

import (
	"encoding/json"
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/ipfs/go-log"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Handler struct {
	pb            *pubsub.PubSub
	oracleTopic   string
	networkTopics mapset.Set
	peerID        peer.ID
	identityMap   map[peer.ID]string
	Logger        *log.ZapEventLogger
	PbMutex       sync.Mutex
}

// TextMessage is more end-user model of regular text messages
type TextMessage struct {
	Topic string  `json:"topic"`
	Body  string  `json:"body"`
	From  peer.ID `json:"from"`
}

func NewHandler(pb *pubsub.PubSub, oracleTopic string, peerID peer.ID, networkTopics mapset.Set) *Handler {
	handler := &Handler{
		pb:            pb,
		oracleTopic:   oracleTopic,
		networkTopics: networkTopics,
		peerID:        peerID,
		identityMap:   make(map[peer.ID]string),
		Logger:        log.Logger("rendezvous"),
	}
	return handler
}

func (h *Handler) HandleIncomingMessage(topic string, msg pubsub.Message, handleTextMessage func(TextMessage)) {
	fromPeerID, err := peer.IDFromBytes(msg.From)
	if err != nil {
		h.Logger.Warn("Error occurred when reading message from field...")
		return
	}
	message := &BaseMessage{}
	if err = json.Unmarshal(msg.Data, message); err != nil {
		h.Logger.Warn("Error occurred during unmarshalling the base message data")
		return
	}
	if message.To != "" && message.To != h.peerID {
		return // Drop message, because it is not for us
	}

	switch message.Flag {
	// Getting regular message
	case FlagGenericMessage:
		textMessage := TextMessage{
			Topic: topic,
			Body:  message.Body,
			From:  fromPeerID,
		}
		handleTextMessage(textMessage)
	// Getting topic request, answer topic response
	case FlagTopicsRequest:
		respond := &GetTopicsRespondMessage{
			BaseMessage: BaseMessage{
				Body: "",
				Flag: FlagTopicsResponse,
				To:   fromPeerID,
			},
			Topics: h.GetTopics(),
		}
		sendData, err := json.Marshal(respond)
		if err != nil {
			h.Logger.Warn("Error occurred during marshalling the respond from TopicsRequest")
			return
		}
		go func() {
			h.PbMutex.Lock()
			if err = h.pb.Publish(h.oracleTopic, sendData); err != nil {
				h.Logger.Warn("Failed to send new message to pubsub topic", err)
			}
			h.PbMutex.Unlock()
		}()
	// Getting topic respond, adding topics to `networkTopics`
	case FlagTopicsResponse:
		respond := &GetTopicsRespondMessage{}
		if err = json.Unmarshal(msg.Data, respond); err != nil {
			h.Logger.Warn("Error occurred during unmarshalling the message data from TopicsResponse")
			return
		}
		for i := 0; i < len(respond.Topics); i++ {
			h.networkTopics.Add(respond.Topics[i])
		}
	// Getting identity request, answer identity response
	case FlagIdentityRequest:
		h.sendIdentityResponse(h.oracleTopic, fromPeerID)
	// Getting identity respond, mapping Multiaddress/MatrixID
	case FlagIdentityResponse:
		h.identityMap[peer.ID(fromPeerID.String())] = message.From.String()
	case FlagGreeting:
		h.Logger.Info("Greetings from " + fromPeerID.String() + " in topic " + topic)
		h.sendIdentityResponse(topic, fromPeerID)
	case FlagGreetingRespond:
		h.Logger.Info("Greeting respond from " + fromPeerID.String() + ":" + message.From.String() + " in topic " + topic)
	case FlagFarewell:
		h.Logger.Info("Greeting respond from " + fromPeerID.String() + ":" + message.From.String() + " in topic " + topic)
	default:
		h.Logger.Info("\nUnknown message type: %#x\n", message.Flag)
	}
}

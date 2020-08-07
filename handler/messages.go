package handler

import (
	"encoding/json"

	"github.com/Secured-Finance/p2p-oracle-node/rpcclient"
	"github.com/libp2p/go-libp2p-core/peer"
)

const (
	FlagGenericMessage   int = 0x0
	FlagTopicsRequest    int = 0x1
	FlagTopicsResponse   int = 0x2
	FlagIdentityRequest  int = 0x3
	FlagIdentityResponse int = 0x4
	FlagGreeting         int = 0x5
	FlagFarewell         int = 0x6
	FlagGreetingRespond  int = 0x7
	FlagEventMessage     int = 0x8
)

// BaseMessage is the basic message format of our protocol
type BaseMessage struct {
	Body *rpcclient.OracleEvent `json:"body"`
	To   peer.ID                `json:"to"`
	Flag int                    `json:"flag"`
	From peer.ID                `json:"from"`
}

// GetTopicsRespondMessage is the format of the message to answer of request for topics
// Flag: 0x2
type GetTopicsRespondMessage struct {
	BaseMessage
	Topics []string `json:"topics"`
}

func (h *Handler) sendIdentityResponse(topic string, fromPeerID peer.ID) {
	var flag int
	if topic == h.oracleTopic {
		flag = FlagIdentityResponse
	} else {
		flag = FlagGreetingRespond
	}
	respond := &BaseMessage{
		Body: &rpcclient.OracleEvent{},
		Flag: flag,
		From: "",
		To:   fromPeerID,
	}
	sendData, err := json.Marshal(respond)
	if err != nil {
		h.Logger.Warn("Error occurred during marshalling the respond from IdentityRequest")
		return
	}
	go func() {
		h.PbMutex.Lock()
		if err = h.pb.Publish(topic, sendData); err != nil {
			h.Logger.Warn("Failed to send new message to pubsub topic", err)
		}
		h.PbMutex.Unlock()
	}()
}

// Requests MatrixID from specific peer
// TODO: refactor with promise
func (h *Handler) RequestPeerIdentity(peerID peer.ID) {
	requestPeersIdentity := &BaseMessage{
		Body: &rpcclient.OracleEvent{},
		To:   peerID,
		Flag: FlagIdentityRequest,
		From: h.peerID,
	}

	h.SendMessageToServiceTopic(requestPeersIdentity)
}

// TODO: refactor
func (h *Handler) SendGreetingInTopic(topic string) {
	greetingMessage := &BaseMessage{
		Body: &rpcclient.OracleEvent{},
		To:   "",
		Flag: FlagGreeting,
		From: h.peerID,
	}

	h.SendMessageToTopic(topic, greetingMessage)
}

// TODO: refactor
func (h *Handler) SendFarewellInTopic(topic string) {
	farewellMessage := &BaseMessage{
		Body: &rpcclient.OracleEvent{},
		To:   "",
		Flag: FlagFarewell,
		From: h.peerID,
	}

	h.SendMessageToTopic(topic, farewellMessage)
}

// Sends marshaled message to the service topic
func (h *Handler) SendMessageToServiceTopic(message *BaseMessage) {
	h.SendMessageToTopic(h.oracleTopic, message)
}

func (h *Handler) SendMessageToTopic(topic string, message *BaseMessage) {
	sendData, err := json.Marshal(message)
	if err != nil {
		h.Logger.Warn("Failed to send message to topic", err)
		return
	}

	go func() {
		h.PbMutex.Lock()
		if err = h.pb.Publish(topic, sendData); err != nil {
			h.Logger.Warn("Failed to send new message to pubsub topic", err)
		}
		h.PbMutex.Unlock()
	}()
}

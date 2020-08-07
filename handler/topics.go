package handler

import "github.com/Secured-Finance/p2p-oracle-node/rpcclient"

// Get list of topics **this** node is subscribed to
func (h *Handler) GetTopics() []string {
	topics := h.pb.GetTopics()
	return topics
}

// Requesting topics from **other** peers
func (h *Handler) RequestNetworkTopics() {
	requestTopicsMessage := &BaseMessage{
		Body: &rpcclient.OracleEvent{},
		Flag: FlagTopicsRequest,
		To:   "",
		From: h.peerID,
	}

	h.SendMessageToServiceTopic(requestTopicsMessage)
}

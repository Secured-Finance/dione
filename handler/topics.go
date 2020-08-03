package handler

// Get list of topics **this** node is subscribed to
func (h *Handler) GetTopics() []string {
	topics := h.pb.GetTopics()
	return topics
}

// Requesting topics from **other** peers
func (h *Handler) RequestNetworkTopics() {
	requestTopicsMessage := &BaseMessage{
		Body: "",
		Flag: FlagTopicsRequest,
		To:   "",
		From: h.peerID,
	}

	h.sendMessageToServiceTopic(requestTopicsMessage)
}

package handler

import peer "github.com/libp2p/go-libp2p-core/peer"

// Returns copy of handler's identity map ([peer.ID]=>[matrixID])
func (h *Handler) GetIdentityMap() map[peer.ID]string {
	return h.identityMap
}

// Get list of peers subscribed on specific topic
func (h *Handler) GetPeers(topic string) []peer.ID {
	peers := h.pb.ListPeers(topic)
	return peers
}

// Blacklists a peer by its id
func (h *Handler) BlacklistPeer(pid peer.ID) {
	h.pb.BlacklistPeer(pid)
}

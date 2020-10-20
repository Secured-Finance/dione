package node

import "github.com/Secured-Finance/p2p-oracle-node/models"

type Handler interface {
	HandleMessage(message *models.Message)
}

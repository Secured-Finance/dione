package pubsub

import (
	"github.com/Secured-Finance/dione/consensus/types"
)

type Handler func(message *types.Message)

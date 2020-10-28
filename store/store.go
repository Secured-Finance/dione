package store

import (
	"database/sql"

	"github.com/Secured-Finance/dione/node"
)

type Store struct {
	db        *sql.DB
	node      *node.Node
	genesisTs uint64
	// genesisTask *types.DioneTask
}

func NewStore(db *sql.DB, node *node.Node, genesisTs uint64) *Store {
	return &Store{
		db:        db,
		node:      node,
		genesisTs: genesisTs,
	}
}

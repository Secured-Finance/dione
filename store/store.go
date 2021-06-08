package store

//import (
//	"github.com/Secured-Finance/dione/node"
//	"github.com/jmoiron/sqlx"
//	_ "github.com/mattn/go-sqlite3"
//)
//
//type Store struct {
//	db        *sqlx.DB
//	node      *node.Node
//	genesisTs uint64
//}
//
//func NewStore(node *node.Node, genesisTs uint64) (*Store, error) {
//	db, err := newDB(node.Config.Store.DatabaseURL)
//	if err != nil {
//		return nil, err
//	}
//
//	defer db.Close()
//
//	return &Store{
//		db:        db,
//		node:      node,
//		genesisTs: genesisTs,
//	}, nil
//}
//
//func newDB(databaseURL string) (*sqlx.DB, error) {
//	db, err := sqlx.Connect("sqlite3", databaseURL)
//	if err != nil {
//		return nil, err
//	}
//
//	if err := db.Ping(); err != nil {
//		return nil, err
//	}
//
//	return db, nil
//}
//
//	TODO: Discuss with ChronosX88 about using custom database to decrease I/O bound
//	specify the migrations for stake storage;

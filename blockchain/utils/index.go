package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/ledgerwatch/lmdb-go/lmdb"
)

const (
	DefaultIndexPrefix = "indexes/"
)

var (
	ErrIndexKeyNotFound = fmt.Errorf("key is not found in the index")
)

type Index struct {
	name  string
	dbEnv *lmdb.Env
	db    lmdb.DBI
}

func NewIndex(name string, dbEnv *lmdb.Env, db lmdb.DBI) *Index {
	return &Index{
		name:  name,
		db:    db,
		dbEnv: dbEnv,
	}
}

func (i *Index) PutUint64(key []byte, value uint64) error {
	return i.dbEnv.Update(func(txn *lmdb.Txn) error {
		var data []byte
		binary.LittleEndian.PutUint64(data, value)
		return txn.Put(i.db, i.constructIndexKey(key), data, 0)
	})
}

func (i *Index) GetUint64(key []byte) (uint64, error) {
	var num uint64
	err := i.dbEnv.View(func(txn *lmdb.Txn) error {
		data, err := txn.Get(i.db, i.constructIndexKey(key))
		if err != nil {
			if lmdb.IsNotFound(err) {
				return ErrIndexKeyNotFound
			}
			return err
		}
		num = binary.LittleEndian.Uint64(data)
		return nil
	})
	if err != nil {
		return 0, err
	}

	return num, nil
}

func (i *Index) PutBytes(key []byte, value []byte) error {
	return i.dbEnv.Update(func(txn *lmdb.Txn) error {
		return txn.Put(i.db, i.constructIndexKey(key), value, 0)
	})
}

func (i *Index) GetBytes(key []byte) ([]byte, error) {
	var data []byte
	err := i.dbEnv.View(func(txn *lmdb.Txn) error {
		valueData, err := txn.Get(i.db, i.constructIndexKey(key))
		if err != nil {
			if lmdb.IsNotFound(err) {
				return ErrIndexKeyNotFound
			}
			return err
		}
		data = valueData
		return nil
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (i *Index) Delete(key []byte) error {
	return i.dbEnv.Update(func(txn *lmdb.Txn) error {
		return txn.Del(i.db, i.constructIndexKey(key), nil)
	})
}

func (i *Index) constructIndexKey(key []byte) []byte {
	k := hex.EncodeToString(key)
	return []byte(fmt.Sprintf("%s/%s/%s", DefaultIndexPrefix, i.name, k))
}

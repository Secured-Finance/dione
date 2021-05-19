package pool

import (
	"encoding/hex"
	"sort"
	"sync"
	"time"

	types2 "github.com/Secured-Finance/dione/blockchain/types"

	"github.com/Secured-Finance/dione/consensus/policy"

	"github.com/Secured-Finance/dione/cache"
)

const (
	DefaultTxTTL    = 10 * time.Minute
	DefaultTxPrefix = "tx_"
)

type Mempool struct {
	m             sync.RWMutex
	cache         cache.Cache
	txDescriptors []string // list of txs in cache
}

func NewMempool(c cache.Cache) (*Mempool, error) {
	mp := &Mempool{
		cache: c,
	}

	var txDesc []string
	err := c.Get("tx_list", &txDesc)
	if err != nil || err != cache.ErrNilValue {
		return nil, err
	}
	mp.txDescriptors = txDesc

	return mp, nil
}

func (mp *Mempool) StoreTx(tx *types2.Transaction) error {
	mp.m.Lock()
	defer mp.m.Unlock()

	hashStr := hex.EncodeToString(tx.Hash)
	err := mp.cache.StoreWithTTL(DefaultTxPrefix+hashStr, tx, DefaultTxTTL)
	mp.txDescriptors = append(mp.txDescriptors, hashStr)
	mp.cache.Store("tx_list", mp.txDescriptors) // update tx list in cache
	return err
}

func (mp *Mempool) GetTxsForNewBlock() []*types2.Transaction {
	mp.m.Lock()
	defer mp.m.Unlock()

	var txForBlock []*types2.Transaction
	var allTxs []*types2.Transaction

	for i, v := range mp.txDescriptors {
		var tx types2.Transaction
		err := mp.cache.Get(DefaultTxPrefix+v, &tx)
		if err != nil {
			if err == cache.ErrNilValue {
				// descriptor is broken
				// delete it and update list
				mp.txDescriptors = removeItemFromStringSlice(mp.txDescriptors, i)
				mp.cache.Store("tx_list", mp.txDescriptors) // update tx list in cache
			}
			continue
		}
		allTxs = append(allTxs, &tx)
	}
	sort.Slice(allTxs, func(i, j int) bool {
		return allTxs[i].Timestamp.Before(allTxs[j].Timestamp)
	})

	for i := 0; i < policy.BlockMaxTransactionCount; i++ {
		if len(allTxs) == 0 {
			break
		}
		tx := allTxs[0]     // get oldest tx
		allTxs = allTxs[1:] // pop tx
		txForBlock = append(txForBlock, tx)
	}

	return txForBlock
}

func removeItemFromStringSlice(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

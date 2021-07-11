package drand

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/Arceliar/phony"

	"github.com/Secured-Finance/dione/beacon"
	"github.com/drand/drand/chain"
	"github.com/drand/drand/client"
	httpClient "github.com/drand/drand/client/http"
	libp2pClient "github.com/drand/drand/lp2p/client"
	"github.com/drand/kyber"
	logging "github.com/ipfs/go-log"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"

	dlog "github.com/drand/drand/log"
	kzap "github.com/go-kit/kit/log/zap"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/lib"
	types "github.com/Secured-Finance/dione/types"
)

var log = logrus.WithFields(logrus.Fields{
	"subsystem": "drand",
})

type DrandBeacon struct {
	phony.Inbox
	DrandClient        client.Client
	PublicKey          kyber.Point
	drandResultChannel <-chan client.Result
	beaconEntryChannel chan types.BeaconEntry
	cacheLock          sync.Mutex
	localCache         map[uint64]types.BeaconEntry
	latestDrandRound   uint64
}

func NewDrandBeacon(ps *pubsub.PubSub) (*DrandBeacon, error) {
	cfg := config.NewDrandConfig()

	drandChain, err := chain.InfoFromJSON(bytes.NewReader([]byte(cfg.ChainInfo)))
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal drand chain info: %w", err)
	}

	dlogger := dlog.NewKitLoggerFrom(kzap.NewZapSugarLogger(
		logging.Logger("drand").SugaredLogger.Desugar(), zapcore.InfoLevel))

	var clients []client.Client
	for _, url := range cfg.Servers {
		client, err := httpClient.NewWithInfo(url, drandChain, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create http drand client: %w", err)
		}
		clients = append(clients, client)
	}

	opts := []client.Option{
		client.WithChainInfo(drandChain),
		client.WithCacheSize(1024),
		client.WithAutoWatch(),
		client.WithLogger(dlogger),
	}

	if ps != nil {
		opts = append(opts, libp2pClient.WithPubsub(ps))
	} else {
		log.Info("Initiated drand with PubSub")
	}

	drandClient, err := client.Wrap(clients, opts...)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create Drand clients")
	}

	db := &DrandBeacon{
		DrandClient: drandClient,
		localCache:  make(map[uint64]types.BeaconEntry),
	}

	db.PublicKey = drandChain.PublicKey

	db.drandResultChannel = db.DrandClient.Watch(context.TODO())
	db.beaconEntryChannel = make(chan types.BeaconEntry)
	err = db.getLatestDrandResult()
	if err != nil {
		return nil, err
	}
	go db.loop(context.TODO())

	return db, nil
}

func (db *DrandBeacon) getLatestDrandResult() error {
	latestDround, err := db.DrandClient.Get(context.TODO(), 0)
	if err != nil {
		log.Errorf("failed to get latest drand round: %v", err)
		return err
	}
	db.cacheValue(newBeaconEntryFromDrandResult(latestDround))
	db.updateLatestDrandRound(latestDround.Round())
	return nil
}

func (db *DrandBeacon) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case res := <-db.drandResultChannel:
			{
				db.cacheValue(newBeaconEntryFromDrandResult(res))
				db.updateLatestDrandRound(res.Round())
				db.newEntry(res)
			}
		}
	}
}

func (db *DrandBeacon) Entry(ctx context.Context, round uint64) (types.BeaconEntry, error) {
	if round != 0 {
		be := db.getCachedValue(round)
		if be != nil {
			return *be, nil
		}
	}

	start := lib.Clock.Now()
	log.Infof("start fetching randomness: round %v", round)
	resp, err := db.DrandClient.Get(ctx, round)
	if err != nil {
		return types.BeaconEntry{}, fmt.Errorf("drand failed Get request: %w", err)
	}
	log.Infof("done fetching randomness: round %v, took %v", round, lib.Clock.Since(start))
	return newBeaconEntryFromDrandResult(resp), nil
}
func (db *DrandBeacon) cacheValue(res types.BeaconEntry) {
	db.cacheLock.Lock()
	defer db.cacheLock.Unlock()
	db.localCache[res.Round] = res
}

func (db *DrandBeacon) getCachedValue(round uint64) *types.BeaconEntry {
	db.cacheLock.Lock()
	defer db.cacheLock.Unlock()
	v, ok := db.localCache[round]
	if !ok {
		return nil
	}
	return &v
}

func (db *DrandBeacon) updateLatestDrandRound(round uint64) {
	db.cacheLock.Lock()
	defer db.cacheLock.Unlock()
	db.latestDrandRound = round
}

func (db *DrandBeacon) VerifyEntry(curr, prev types.BeaconEntry) error {
	if prev.Round == 0 {
		return nil
	}
	if be := db.getCachedValue(curr.Round); be != nil {
		return nil
	}
	b := &chain.Beacon{
		PreviousSig: prev.Metadata["signature"].([]byte),
		Round:       curr.Round,
		Signature:   curr.Metadata["signature"].([]byte),
	}
	return chain.VerifyBeacon(db.PublicKey, b)
}

func (db *DrandBeacon) LatestBeaconRound() uint64 {
	db.cacheLock.Lock()
	defer db.cacheLock.Unlock()
	return db.latestDrandRound
}

func (db *DrandBeacon) newEntry(res client.Result) {
	db.Act(nil, func() {
		db.beaconEntryChannel <- types.NewBeaconEntry(res.Round(), res.Randomness(), map[string]interface{}{"signature": res.Signature()})
	})
}

func (db *DrandBeacon) NewEntries() <-chan types.BeaconEntry {
	return db.beaconEntryChannel
}

func newBeaconEntryFromDrandResult(res client.Result) types.BeaconEntry {
	return types.NewBeaconEntry(res.Round(), res.Randomness(), map[string]interface{}{"signature": res.Signature()})
}

var _ beacon.BeaconAPI = (*DrandBeacon)(nil)

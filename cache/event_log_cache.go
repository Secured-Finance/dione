package cache

import (
	oracleEmitter "github.com/Secured-Finance/dione/contracts/oracleemitter"
	"github.com/VictoriaMetrics/fastcache"
	"github.com/fxamacker/cbor/v2"
)

const (
	// in megabytes
	DefaultEventLogCacheCapacity = 32000000
)

type EventLogCache struct {
	cache *fastcache.Cache
}

func NewEventLogCache() *EventLogCache {
	return &EventLogCache{
		cache: fastcache.New(DefaultEventLogCacheCapacity),
	}
}

func (elc *EventLogCache) Store(key string, event interface{}) error {
	mRes, err := cbor.Marshal(event)
	if err != nil {
		return err
	}

	elc.cache.SetBig([]byte(key), mRes)

	return nil
}

func (elc *EventLogCache) GetOracleRequestEvent(key string) (*oracleEmitter.OracleEmitterNewOracleRequest, error) {
	var mData []byte
	mData = elc.cache.GetBig(mData, []byte(key))

	var event *oracleEmitter.OracleEmitterNewOracleRequest
	err := cbor.Unmarshal(mData, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (elc *EventLogCache) Delete(key string) {
	elc.cache.Del([]byte(key))
}
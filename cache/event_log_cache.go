package cache

import (
	oracleEmitter "github.com/Secured-Finance/dione/contracts/oracleemitter"
	"github.com/VictoriaMetrics/fastcache"
	"github.com/fxamacker/cbor/v2"
	"github.com/go-redis/redis/v8"
)

const (
	// in megabytes
	DefaultEventLogCacheCapacity = 32000000
)
const (
	RedisAddr = "redisDB:6379"
)


type EventLogCache struct {
	cache *fastcache.Cache
}

type EventRedisCache struct {
	Client *redis.Client
}

func NewEventLogCache() *EventLogCache {
	return &EventLogCache{
		cache: fastcache.New(DefaultEventLogCacheCapacity),
	}
}

func NewEventRedisCache() *EventRedisCache {
	client := redis.NewClient(&redis.Options{
	   Addr: RedisAddr,
	   Password: "",
	   DB: 0,
	})

	return &EventRedisCache{
	   Client: client,
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

func (erc *EventLogCache) GetOracleRequestEvent(key string) (*oracleEmitter.OracleEmitterNewOracleRequest, error) {
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

func (erc *EventRedisCache) StoreRedis(key string, event interface{}) error {
	mRes, err := cbor.Marshal(event)
	if err != nil {
		return err
	}

	erc.Client.Set(key, mRes)

	return nil
}

func (elc *EventRedisCache) GetOracleRequestEventFromRedis(key string) (*oracleEmitter.OracleEmitterNewOracleRequest, error) {
	var mData []byte
	mData = erc.Client.Get(key)

	var event *oracleEmitter.OracleEmitterNewOracleRequest
	err := cbor.Unmarshal(mData, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (elc *EventRedisCache) DeleteRedisKey(key string) {
	erc.Client.Del([]byte(key))
}
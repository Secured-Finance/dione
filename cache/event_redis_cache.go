package cache

import (
	oracleEmitter "github.com/Secured-Finance/dione/contracts/oracleemitter"
	"github.com/Secured-Finance/dione/config"
	"github.com/fxamacker/cbor/v2"
	"github.com/go-redis/redis/v8"
)

type EventRedisCache struct {
	Client *redis.Client
}

func NewEventRedisCache(config *config.Config) *EventRedisCache {
	client := redis.NewClient(&redis.Options{
	   Addr: config.Redis.Addr,
	   Password: config.Redis.Password,
	   DB: config.Redis.DB,
	})

	return &EventRedisCache{
	   Client: client,
	}
}

 func (erc *EventRedisCache) Store(key string, event interface{}) error {
	mRes, err := cbor.Marshal(event)
	if err != nil {
		return err
	}

	erc.Client.Set(key, mRes)

	return nil
}

func (erc *EventRedisCache) GetOracleRequestEvent(key string) (*oracleEmitter.OracleEmitterNewOracleRequest, error) {
	var mData []byte
	mData = erc.Client.Get(key)

	var event *oracleEmitter.OracleEmitterNewOracleRequest
	err := cbor.Unmarshal(mData, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (erc *EventRedisCache) Delete(key string) {
	erc.Client.Del(key)
}
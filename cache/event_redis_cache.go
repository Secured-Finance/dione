package cache

import (
	"context"

	"github.com/Secured-Finance/dione/config"
	"github.com/Secured-Finance/dione/contracts/dioneOracle"
	"github.com/fxamacker/cbor/v2"
	"github.com/go-redis/redis/v8"
)

type EventRedisCache struct {
	Client *redis.Client
	ctx    context.Context
}

func NewEventRedisCache(config *config.Config) *EventRedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	return &EventRedisCache{
		Client: client,
		ctx:    context.Background(),
	}
}

func (erc *EventRedisCache) Store(key string, event interface{}) error {
	mRes, err := cbor.Marshal(event)
	if err != nil {
		return err
	}

	erc.Client.Set(erc.ctx, key, mRes, 0)

	return nil
}

func (erc *EventRedisCache) GetOracleRequestEvent(key string) (*dioneOracle.DioneOracleNewOracleRequest, error) {
	mData, err := erc.Client.Get(erc.ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var event *dioneOracle.DioneOracleNewOracleRequest
	err = cbor.Unmarshal(mData, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (erc *EventRedisCache) Delete(key string) {
	erc.Client.Del(erc.ctx, key)
}

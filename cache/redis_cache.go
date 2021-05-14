package cache

import (
	"context"
	"time"

	"github.com/Secured-Finance/dione/config"
	"github.com/fxamacker/cbor/v2"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	Client *redis.Client
	ctx    context.Context
}

func NewRedisCache(config *config.Config) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	return &RedisCache{
		Client: client,
		ctx:    context.Background(),
	}
}

func (rc *RedisCache) Store(key string, value interface{}) error {
	mRes, err := cbor.Marshal(value)
	if err != nil {
		return err
	}

	rc.Client.Set(rc.ctx, key, mRes, 0)

	return nil
}

func (rc *RedisCache) StoreWithTTL(key string, value interface{}, ttl time.Duration) error {
	mRes, err := cbor.Marshal(value)
	if err != nil {
		return err
	}

	rc.Client.Set(rc.ctx, key, mRes, ttl)

	return nil
}

func (rc *RedisCache) Get(key string, value interface{}) error {
	data, err := rc.Client.Get(rc.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrNilValue
		}
		return err
	}
	err = cbor.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) Delete(key string) {
	rc.Client.Del(rc.ctx, key)
}

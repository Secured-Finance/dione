package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	"github.com/Secured-Finance/dione/config"
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
	data, err := gobMarshal(value)
	if err != nil {
		return err
	}

	rc.Client.Set(rc.ctx, key, data, 0)

	return nil
}

func (rc *RedisCache) StoreWithTTL(key string, value interface{}, ttl time.Duration) error {
	data, err := gobMarshal(value)
	if err != nil {
		return err
	}
	rc.Client.Set(rc.ctx, key, data, ttl)

	return nil
}

func gobMarshal(val interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(val)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gobUnmarshal(data []byte, val interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(&val)
}

func (rc *RedisCache) Get(key string, value interface{}) error {
	data, err := rc.Client.Get(rc.ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrNotFound
		}
		return err
	}

	return gobUnmarshal(data, &value)
}

func (rc *RedisCache) Delete(key string) {
	rc.Client.Del(rc.ctx, key)
}

func (rc *RedisCache) Items() map[string]interface{} {
	return nil // TODO
}

package cache

import (
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/fxamacker/cbor/v2"
)

const (
	// DefaultInMemoryCacheCapacity is maximal in-memory cache size in bytes
	DefaultInMemoryCacheCapacity = 32000000
)

type InMemoryCache struct {
	cache *fastcache.Cache
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: fastcache.New(DefaultInMemoryCacheCapacity),
	}
}

func (imc *InMemoryCache) Store(key string, value interface{}) error {
	mRes, err := cbor.Marshal(value)
	if err != nil {
		return err
	}

	imc.cache.SetBig([]byte(key), mRes)

	return nil
}

func (imc *InMemoryCache) StoreWithTTL(key string, value interface{}, ttl time.Duration) error {
	return imc.Store(key, value) // fastcache doesn't support ttl for values
}

func (imc *InMemoryCache) Get(key string, v interface{}) error {
	data := make([]byte, 0)
	imc.cache.GetBig(data, []byte(key))
	if len(data) == 0 {
		return ErrNilValue
	}
	err := cbor.Unmarshal(data, v)
	if err != nil {
		return err
	}

	return nil
}

func (imc *InMemoryCache) Delete(key string) {
	imc.cache.Del([]byte(key))
}

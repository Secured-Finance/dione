package cache

import (
	"fmt"
	"reflect"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	DefaultCacheExpiration = 5 * time.Minute
	DefaultGCInterval      = 10 * time.Minute
)

type InMemoryCache struct {
	cache *cache.Cache
}

func NewInMemoryCache() Cache {
	return &InMemoryCache{
		cache: cache.New(DefaultCacheExpiration, DefaultGCInterval),
	}
}

func (imc *InMemoryCache) Store(key string, value interface{}) error {
	imc.cache.Set(key, value, cache.NoExpiration)

	return nil
}

func (imc *InMemoryCache) StoreWithTTL(key string, value interface{}, ttl time.Duration) error {
	imc.cache.Set(key, value, ttl)
	return nil
}

func (imc *InMemoryCache) Get(key string, value interface{}) error {
	v, exists := imc.cache.Get(key)
	if !exists {
		return ErrNotFound
	}
	reflectedValue := reflect.ValueOf(value)
	if reflectedValue.Kind() != reflect.Ptr {
		return fmt.Errorf("value isn't a pointer")
	}
	if reflectedValue.IsNil() {
		reflectedValue.Set(reflect.New(reflectedValue.Type().Elem()))
	}
	reflectedValue.Elem().Set(reflect.ValueOf(v).Elem())

	return nil
}

func (imc *InMemoryCache) Delete(key string) {
	imc.cache.Delete(key)
}

func (imc *InMemoryCache) Items() map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range imc.cache.Items() {
		m[k] = v.Object
	}
	return m
}

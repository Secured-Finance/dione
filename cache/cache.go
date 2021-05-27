package cache

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("key doesn't exist in cache")

type Cache interface {
	Store(key string, value interface{}) error
	StoreWithTTL(key string, value interface{}, ttl time.Duration) error
	Get(key string, value interface{}) error
	Delete(key string)
}

package cache

import (
	"context"
	"time"
)

const (
	INMEMORYCACHE = "m"
	REDISCACHE    = "r"
)

type Cache interface {
	Get(context.Context, string) (any, error)
	Set(context.Context, string, any, time.Duration) error
	Delete(context.Context, string) error
	SetTTL(context.Context, string, time.Duration) error
	GetTTL(context.Context, string) (time.Duration, error)
	Exists(context.Context, string) (bool, error)
	Clear(context.Context) error
	Close() error
	Description() string
}

func NewCache(cacheType string) (Cache, error) {
	switch cacheType {
	case INMEMORYCACHE:
		return NewInMemoryCache(), nil
	case REDISCACHE:
		return NewRedisCache("", "", 0)
	default:
		return NewDummyCache(), nil
	}
}

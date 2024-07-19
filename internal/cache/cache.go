package cache

import (
	"context"
	"time"
)

const (
	INMEMORYCACHE = "memory"
	REDISCACHE    = "redis"
	DUMMYCACHE    = "dummy"
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

func NewCache(cacheType string, addr, password string, db int) (Cache, error) {
	switch cacheType {
	case INMEMORYCACHE:
		return NewInMemoryCache(), nil
	case REDISCACHE:
		return NewRedisCache(addr, password, db)
	default:
		return NewDummyCache(), nil
	}
}

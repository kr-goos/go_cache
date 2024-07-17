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
	Exists(string) (bool, error)
	Clear() error
	Close()
	Description() string
}

func NewCache(cacheType string) Cache {
	switch cacheType {
	case INMEMORYCACHE:
		return NewInMemoryCache()
	case REDISCACHE:
		return nil
	default:
		return NewDummyCache()
	}
}

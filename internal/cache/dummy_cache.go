package cache

import (
	"context"
	"time"
)

type dummyCache struct{}

func NewDummyCache() *dummyCache {
	return &dummyCache{}
}

func (c *dummyCache) Get(ctx context.Context, key string) (any, error) {
	return nil, nil
}

func (c *dummyCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return nil
}

func (c *dummyCache) Delete(ctx context.Context, key string) error {
	return nil
}

func (c *dummyCache) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	return nil
}

func (c *dummyCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, nil
}

func (c *dummyCache) Exists(key string) (bool, error) {
	return false, nil
}

func (c *dummyCache) Clear() error {
	return nil
}

func (c *dummyCache) Close() {}

func (c *dummyCache) Description() string {
	return "DummyCache: A dummy cache implementation that does nothing"
}

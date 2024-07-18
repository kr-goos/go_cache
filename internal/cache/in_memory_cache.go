package cache

import (
	"context"
	"sync"
	"time"
)

type cacheItem struct {
	value  any
	expiry time.Time
}

type InMemoryCache struct {
	items map[string]*cacheItem
	mtx   sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		items: make(map[string]*cacheItem),
	}
}

func (c *InMemoryCache) Get(_ context.Context, key string) (any, error) {
	c.mtx.RLock()
	item, exists := c.items[key]
	if !exists {
		c.mtx.RUnlock()
		return nil, ErrKeyNotFound
	}

	if item.expiry.Before(time.Now()) {
		c.mtx.RUnlock()
		c.mtx.Lock()
		delete(c.items, key)
		c.mtx.Unlock()
		return nil, ErrKeyExpired
	}
	c.mtx.RUnlock()
	return item.value, nil
}

func (c *InMemoryCache) Set(_ context.Context, key string, value any, ttl time.Duration) error {

	expiry := time.Now().Add(ttl)
	item := &cacheItem{
		value:  value,
		expiry: expiry,
	}
	c.mtx.Lock()
	c.items[key] = item
	c.mtx.Unlock()
	return nil
}

func (c *InMemoryCache) Delete(_ context.Context, key string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.items, key)
	return nil
}

func (c *InMemoryCache) SetTTL(_ context.Context, key string, ttl time.Duration) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	item, exists := c.items[key]
	if !exists {
		return ErrKeyNotFound
	}
	item.expiry = time.Now().Add(ttl)
	return nil
}

func (c *InMemoryCache) GetTTL(_ context.Context, key string) (time.Duration, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return 0, ErrKeyNotFound
	}
	ttl := time.Until(item.expiry)
	return ttl, nil
}

func (c *InMemoryCache) Exists(_ context.Context, key string) (bool, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	_, exists := c.items[key]
	return exists, nil
}

func (c *InMemoryCache) Clear(_ context.Context) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.items = make(map[string]*cacheItem)
	return nil
}

func (c *InMemoryCache) Close() error { return nil }

func (c *InMemoryCache) Description() string {
	return "InMemoryCache: A simple in-memory cache implementation"
}

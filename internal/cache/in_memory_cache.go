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

type inMemoryCache struct {
	items map[string]*cacheItem
	mtx   sync.RWMutex
}

func NewInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		items: make(map[string]*cacheItem),
	}
}

func (c *inMemoryCache) get(key string) (any, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, ErrKeyNotFound
	}
	if item.expiry.Before(time.Now()) {
		delete(c.items, key)
		return nil, ErrKeyExpired
	}
	return item.value, nil
}

func (c *inMemoryCache) set(key string, value any, ttl time.Duration) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	expiry := time.Now().Add(ttl)
	item := &cacheItem{
		value:  value,
		expiry: expiry,
	}
	c.items[key] = item
	return nil
}

func (c *inMemoryCache) delete(key string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.items, key)
	return nil
}

func (c *inMemoryCache) setTTL(key string, ttl time.Duration) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	item, exists := c.items[key]
	if !exists {
		return ErrKeyNotFound
	}
	item.expiry = time.Now().Add(ttl)
	return nil
}

func (c *inMemoryCache) getTTL(key string) (time.Duration, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return 0, ErrKeyNotFound
	}
	ttl := time.Until(item.expiry)
	return ttl, nil
}

func (c *inMemoryCache) Get(ctx context.Context, key string) (any, error) {
	done := make(chan struct{})
	var result any
	var err error

	go func() {
		result, err = c.get(key)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return result, err
	}
}

func (c *inMemoryCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	done := make(chan struct{})
	var err error

	go func() {
		err = c.set(key, value, ttl)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return err
	}
}

func (c *inMemoryCache) Delete(ctx context.Context, key string) error {
	done := make(chan struct{})
	var err error

	go func() {
		err = c.delete(key)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return err
	}
}

func (c *inMemoryCache) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	done := make(chan struct{})
	var err error

	go func() {
		err = c.setTTL(key, ttl)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return err
	}
}

func (c *inMemoryCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	done := make(chan struct{})
	var ttl time.Duration
	var err error

	go func() {
		ttl, err = c.getTTL(key)
		close(done)
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-done:
		return ttl, err
	}
}

func (c *inMemoryCache) Exists(key string) (bool, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	_, exists := c.items[key]
	return exists, nil
}

func (c *inMemoryCache) Clear() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.items = make(map[string]*cacheItem)
	return nil
}

func (c *inMemoryCache) Close() {
	// c.mtx.Lock()
	// defer c.mtx.Unlock()

	// c.items = make(map[string]*cacheItem)
}

func (c *inMemoryCache) Description() string {
	return "InMemoryCache: A simple in-memory cache implementation"
}

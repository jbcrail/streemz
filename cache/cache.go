package cache

import (
	"time"
)

type Cache interface {
	Get(key string) interface{}
	Set(key string, value interface{}, expire time.Duration)
	Exists(key string) bool
	Clear(key string)
	ClearAll()
}

type InMemoryCache struct {
	storage    map[string]interface{}
	expiration map[string]time.Time
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		storage:    make(map[string]interface{}),
		expiration: make(map[string]time.Time),
	}
}

func (c *InMemoryCache) Get(key string) interface{} {
	v, ok := c.storage[key]
	if !ok {
		return nil
	}
	expiration, _ := c.expiration[key]
	if time.Now().After(expiration) {
		return nil
	}
	return v
}

func (c *InMemoryCache) Set(key string, value interface{}, expire time.Duration) {
	expiration := time.Now().Add(expire)
	c.storage[key] = value
	c.expiration[key] = expiration
}

func (c *InMemoryCache) Exists(key string) bool {
	_, ok := c.storage[key]
	return ok
}

func (c *InMemoryCache) Clear(key string) {
	delete(c.storage, key)
	delete(c.expiration, key)
}

func (c *InMemoryCache) ClearAll() {
	c.storage = make(map[string]interface{})
	c.expiration = make(map[string]time.Time)
}

package secretsmanager

import (
	"sync"
	"time"
)

type CachedSecret struct {
	Value     string
	Timestamp time.Time
}

type SecretCache struct {
	cache map[string]CachedSecret
	ttl   time.Duration
	mu    sync.RWMutex
}

func NewSecretCache(ttlSeconds int64) *SecretCache {
	return &SecretCache{
		cache: make(map[string]CachedSecret),
		ttl:   time.Duration(ttlSeconds) * time.Second,
	}
}

func (c *SecretCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, found := c.cache[key]
	if !found {
		return "", false
	}

	if time.Since(cached.Timestamp) > c.ttl {
		return "", false
	}

	return cached.Value, true
}

func (c *SecretCache) Set(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = CachedSecret{
		Value:     value,
		Timestamp: time.Now(),
	}
}

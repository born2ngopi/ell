package memcache

import "sync"

type cache struct {
	token map[string]string
	mu    sync.RWMutex
}

var c *cache

func Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token[key] = value
}

func Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token[key]
}

func Update(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token[key] = value
}

func Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.token, key)
}

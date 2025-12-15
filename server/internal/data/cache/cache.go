package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	items map[string]CacheItem
	mutex sync.RWMutex
}

func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]CacheItem),
	}
	
	// Limpeza automática em segundo plano
	go cache.cleanup()
	
	return cache
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	expiration := time.Now().Add(duration).UnixNano()
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, found := c.items[key]
	if !found {
		return nil, false
	}
	
	// Verifica se expirou
	if time.Now().UnixNano() > item.Expiration {
		return nil, false
	}
	
	return item.Value, true
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.items, key)
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute) // Limpeza a cada 5 minutos
	defer ticker.Stop()
	
	for range ticker.C {
		now := time.Now().UnixNano()
		c.mutex.Lock()
		for k, v := range c.items {
			if now > v.Expiration {
				delete(c.items, k)
			}
		}
		c.mutex.Unlock()
	}
}
package cache

import (
	"sync"
)

type Cache struct {
	store sync.Map
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) Set(Key string, value interface{}) {
	c.store.Store(Key, value)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.store.Load(key)
}

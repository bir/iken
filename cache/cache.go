package cache

import (
	"sync"
)

// Cache is a simple cache and has only supports manual eviction.
type Cache[K comparable, V any] struct {
	items map[K]V //nolint: structcheck
	*sync.RWMutex
}

// NewCache creates a new non-thread safe cache.
func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		items:   make(map[K]V, 0),
		RWMutex: &sync.RWMutex{},
	}
}

// Set sets any item to the cache. replacing any existing item.
// The default item never expires.
func (c *Cache[K, V]) Set(k K, v V) {
	c.Lock()
	defer c.Unlock()

	c.items[k] = v
}

// Get gets an item from the cache.
// Returns the item or zero value, and a bool indicating whether the key was found.
func (c *Cache[K, V]) Get(k K) (V, bool) { //nolint: ireturn
	c.RLock()
	defer c.RUnlock()

	got, found := c.items[k]

	return got, found
}

// Keys returns cache keys. the order is sorted by created.
func (c *Cache[K, _]) Keys() []K {
	c.RLock()
	defer c.RUnlock()

	l := len(c.items)
	if l == 0 {
		return nil
	}

	out := make([]K, 0, l)
	for key := range c.items {
		out = append(out, key)
	}

	return out
}

// Delete deletes the item with provided key from the cache.
func (c *Cache[K, V]) Delete(key K) {
	c.Lock()
	defer c.Unlock()

	delete(c.items, key)
}

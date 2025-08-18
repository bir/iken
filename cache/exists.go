package cache

import (
	"sync"
)

// Exists is a thread safe data structure to track usage of keys.
// The backing store is a map of empty struct{}, this is the simplest memory efficient
// means in Go for tracking keys.
// See: https://medium.com/easyread/golang-series-empty-struct-ed317e6d8600
type Exists[K comparable] struct {
	*sync.RWMutex

	items map[K]struct{}
}

// NewExists creates a new thread safe exists.
func NewExists[K comparable]() *Exists[K] {
	return &Exists[K]{
		items:   make(map[K]struct{}, 0),
		RWMutex: &sync.RWMutex{},
	}
}

// Mark flags the key.
func (c *Exists[K]) Mark(k K) {
	c.Lock()
	defer c.Unlock()

	c.items[k] = struct{}{}
}

// MarkIf flags the key and returns if it was previously set.
func (c *Exists[K]) MarkIf(k K) bool {
	c.Lock()
	defer c.Unlock()

	if _, found := c.items[k]; found {
		return true
	}

	c.items[k] = struct{}{}

	return false
}

// Check returns true if the key is marked.
func (c *Exists[K]) Check(k K) bool {
	c.RLock()
	defer c.RUnlock()

	_, found := c.items[k]

	return found
}

// Keys returns existing keys, the order is indeterminate.
func (c *Exists[K]) Keys() []K {
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
func (c *Exists[K]) Delete(key K) {
	c.Lock()
	defer c.Unlock()

	delete(c.items, key)
}

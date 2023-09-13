package cache

// Cache is the basic contract for all cache implementations.
type Cache[K comparable, V any] interface {
	// Get gets an item from the cache.
	Get(k K) (V, bool)
	// Set sets any item to the cache, replacing any existing item if it exists.
	Set(k K, v V)
	// Delete deletes the item with provided key from the cache.
	Delete(key K)
	// Keys returns existing keys, the order is indeterminate.
	Keys() []K
	// Clear resets the cache.
	Clear()
}

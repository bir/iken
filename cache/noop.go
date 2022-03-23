package cache

// NoOpCache is a facade cache, it never returns a hit.
// Useful for easily disabling cache at runtime.
type NoOpCache[K comparable, V any] struct{}

// NewNoOpCache creates a new NOP cache.
func NewNoOpCache[K comparable, V any]() *NoOpCache[K, V] {
	return &NoOpCache[K, V]{}
}

// Set no-op.
func (c *NoOpCache[K, V]) Set(_ K, _ V) {
}

// Get always returns !ok.
func (c *NoOpCache[K, V]) Get(_ K) (out V, ok bool) { //nolint: ireturn
	return
}

// Keys always returns nil array.
func (c *NoOpCache[K, _]) Keys() []K {
	return nil
}

// Delete no-op.
func (c *NoOpCache[K, V]) Delete(_ K) {
}

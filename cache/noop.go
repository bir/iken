package cache

// NoOp is a facade cache, it never returns a hit.
// Useful for easily disabling cache at runtime.
type NoOp[K comparable, V any] struct{}

// NewNoOp creates a new NOP cache.
func NewNoOp[K comparable, V any]() *NoOp[K, V] {
	return &NoOp[K, V]{}
}

// Set no-op.
func (c *NoOp[K, V]) Set(_ K, _ V) {
}

// Get always returns !ok.
func (c *NoOp[K, V]) Get(_ K) (out V, ok bool) { //nolint: ireturn
	return
}

// Keys always returns nil array.
func (c *NoOp[K, _]) Keys() []K {
	return nil
}

// Delete no-op.
func (c *NoOp[K, V]) Delete(_ K) {
}

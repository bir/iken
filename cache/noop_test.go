package cache_test

import (
	"testing"

	"github.com/bir/iken/cache"
	"github.com/stretchr/testify/assert"
)

// Type assertion
var _ cache.Cache[string, string] = &cache.NoOp[string, string]{}

func TestNoOpCache(t *testing.T) {
	c := cache.NewNoOp[string, int]()
	// Empty
	v, ok := c.Get("a")
	assert.Equal(t, 0, v)
	assert.Equal(t, false, ok)

	kk := c.Keys()
	assert.Equal(t, 0, len(kk))

	// New Value
	c.Set("a", 1)
	v, ok = c.Get("a")
	assert.Equal(t, 0, v)
	assert.Equal(t, false, ok)

	kk = c.Keys()
	assert.Equal(t, 0, len(kk))

	// Override
	c.Set("a", 2)
	v, ok = c.Get("a")
	assert.Equal(t, 0, v)
	assert.Equal(t, false, ok)

	// New Value
	v, ok = c.Get("b")
	assert.Equal(t, 0, v)
	assert.Equal(t, false, ok)

	c.Set("b", 2)
	v, ok = c.Get("b")
	assert.Equal(t, 0, v)
	assert.Equal(t, false, ok)

	kk = c.Keys()
	assert.Equal(t, 0, len(kk))

	// Delete
	c.Delete("a")

	kk = c.Keys()
	assert.Equal(t, 0, len(kk))
}

package cache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/cache"
)

// Type assertion
var _ cache.Cache[string, string] = &cache.NoOp[string, string]{}

func TestNoOpCache(t *testing.T) {
	c := cache.NewNoOp[string, int]()
	// Empty
	v, ok := c.Get("a")
	assert.Zero(t, v)
	assert.False(t, ok)
	kk := c.Keys()
	assert.Empty(t, kk)

	// New Value
	c.Set("a", 1)
	v, ok = c.Get("a")
	assert.Zero(t, v)
	assert.False(t, ok)
	kk = c.Keys()
	assert.Empty(t, kk)

	// Override
	c.Set("a", 2)
	v, ok = c.Get("a")
	assert.Zero(t, v)
	assert.False(t, ok)

	// New Value
	v, ok = c.Get("b")
	assert.Zero(t, v)
	assert.False(t, ok)

	c.Set("b", 2)
	v, ok = c.Get("b")
	assert.Zero(t, v)
	assert.False(t, ok)
	kk = c.Keys()
	assert.Empty(t, kk)

	// Delete
	c.Delete("a")

	kk = c.Keys()
	assert.Empty(t, kk)

	c.Clear()
}

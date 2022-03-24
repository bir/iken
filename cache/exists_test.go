package cache_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/bir/iken/cache"
	"github.com/stretchr/testify/assert"
)

func ExampleExists() {
	c := cache.NewExists[string]()

	// Empty check
	ok := c.Check("a")
	fmt.Println("a exists", ok)

	// Basic Mark and Check
	c.Mark("a")
	ok = c.Check("a")
	fmt.Println("a exists", ok)

	// Add another Mark
	c.Mark("b")
	ok = c.Check("a")
	fmt.Println("a exists", ok)
	ok = c.Check("b")
	fmt.Println("b exists", ok)

	// List keys
	kk := c.Keys()
	sort.Strings(kk)
	fmt.Println("keys", kk)

	// Remove item
	c.Delete("a")
	ok = c.Check("a")
	fmt.Println("a exists", ok)
	ok = c.Check("b")
	fmt.Println("b exists", ok)

	// Output:
	// a exists false
	// a exists true
	// a exists true
	// b exists true
	// keys [a b]
	// a exists false
	// b exists true
}

func TestExists(t *testing.T) {
	c := cache.NewExists[int]()

	assert.False(t, c.Check(1), "empty check")
	assert.Nil(t, c.Keys(), "empty keys")

	c.Mark(99)
	assert.True(t, c.Check(99), "valid check")
	assert.False(t, c.Check(1), "empty check")

	c.Mark(1)
	assert.True(t, c.Check(99), "valid check")
	assert.True(t, c.Check(1), "valid check")

	kk := c.Keys()
	sort.Ints(kk)
	assert.Equal(t, []int{1, 99}, kk, "matching keys")

	c.Delete(1)
	kk = c.Keys()
	sort.Ints(kk)
	assert.True(t, c.Check(99), "valid check")
	assert.False(t, c.Check(1), "empty check")
	assert.Equal(t, []int{99}, kk, "matching keys")
}

package cache_test

import (
	"encoding/base32"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"testing"

	"github.com/bir/iken/cache"
	"github.com/stretchr/testify/assert"
)

func ExampleCache() {
	c := cache.NewBasic[string, int]()
	c.Set("a", 1)
	out, ok := c.Get("a")
	fmt.Println(out, ok)
	out, ok = c.Get("b")
	fmt.Println(out, ok)
	c.Set("b", 2)
	kk := c.Keys()
	fmt.Println(kk)
	c.Delete("a")
	kk = c.Keys()
	fmt.Println(kk)

	// Output:
	// 1 true
	// 0 false
	// [a b]
	// [b]
}

func randString(length int) string {
	randBytes := make([]byte, length)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)
	}

	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randBytes)
}

// Type assertion
var _ cache.Cache[string, string] = cache.NewBasic[string, string]()

func TestCache(t *testing.T) {
	c := cache.NewBasic[string, int]()
	// Empty
	v, ok := c.Get("a")
	assert.Equal(t, 0, v)
	assert.Equal(t, false, ok)

	kk := c.Keys()
	assert.Equal(t, 0, len(kk))

	// New Value
	c.Set("a", 1)
	v, ok = c.Get("a")
	assert.Equal(t, 1, v)
	assert.Equal(t, true, ok)

	kk = c.Keys()
	assert.Equal(t, 1, len(kk))
	assert.Equal(t, []string{"a"}, kk)

	// Override
	c.Set("a", 2)
	v, ok = c.Get("a")
	assert.Equal(t, 2, v)
	assert.Equal(t, true, ok)

	// New Value
	v, ok = c.Get("b")
	assert.Equal(t, 0, v)
	assert.Equal(t, false, ok)

	c.Set("b", 2)
	v, ok = c.Get("b")
	assert.Equal(t, 2, v)
	assert.Equal(t, true, ok)

	kk = c.Keys()
	sort.Strings(kk)
	assert.Equal(t, 2, len(kk))
	assert.Equal(t, []string{"a", "b"}, kk)

	// Delete
	c.Delete("a")

	kk = c.Keys()
	assert.Equal(t, 1, len(kk))
	assert.Equal(t, []string{"b"}, kk)
}

func TestMultiThread(t *testing.T) {
	c := cache.NewBasic[int, string]()
	var wg sync.WaitGroup
	for i := int64(0); i < 100; i++ {
		wg.Add(1)
		go func(i int64) {
			defer wg.Done()
			m := rand.New(rand.NewSource(i))
			for n := 0; n < 10000; n++ {
				key := m.Intn(100)
				value := randString(10)
				c.Set(key, value)
				c.Get(key)
			}
		}(i)
	}

	wg.Wait()
}

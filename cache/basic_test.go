package cache_test

import (
	crand "crypto/rand"
	"encoding/base32"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/cache"
)

func ExampleBasic() {
	c := cache.NewBasic[string, int]()
	c.Set("a", 1)
	out, ok := c.Get("a")
	fmt.Println(out, ok)
	out, ok = c.Get("b")
	fmt.Println(out, ok)
	c.Set("b", 2)
	kk := c.Keys()
	sort.Strings(kk)
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

	_, err := crand.Read(randBytes)
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
	assert.Zero(t, v)
	assert.False(t, ok)

	kk := c.Keys()
	assert.Empty(t, kk)

	// New Value
	c.Set("a", 1)
	v, ok = c.Get("a")
	assert.Equal(t, 1, v)
	assert.True(t, ok)

	kk = c.Keys()
	assert.Len(t, kk, 1)
	assert.Equal(t, []string{"a"}, kk)

	// Override
	c.Set("a", 2)
	v, ok = c.Get("a")
	assert.Equal(t, 2, v)
	assert.True(t, ok)

	// New Value
	v, ok = c.Get("b")
	assert.Zero(t, v)
	assert.False(t, ok)

	c.Set("b", 2)
	v, ok = c.Get("b")
	assert.Equal(t, 2, v)
	assert.True(t, ok)

	kk = c.Keys()
	sort.Strings(kk)
	assert.Len(t, kk, 2)
	assert.Equal(t, []string{"a", "b"}, kk)

	// Delete
	c.Delete("a")

	kk = c.Keys()
	assert.Len(t, kk, 1)
	assert.Equal(t, []string{"b"}, kk)
}

type foo struct {
	c cache.Cache[int, string]
}

func TestMultiThread(t *testing.T) {
	const workers = int64(1000)
	const iterations = 1000
	f := foo{}
	count := atomic.Int64{}
	f.c = cache.NewBasic[int, string]()
	var wg sync.WaitGroup
	for i := range workers {
		wg.Add(1)
		go func(i int64) {
			defer wg.Done()
			f.c.Clear()
			m := rand.New(rand.NewSource(i))
			for range iterations {
				key := m.Intn(100)
				value := randString(10)
				f.c.Set(key, value)
				f.c.Get(key)
				count.Add(1)
			}
		}(i)
	}

	wg.Wait()

	expected := workers * iterations
	if count.Load() != expected {
		t.Errorf("Expected %d, got %d", expected, count.Load())
	}
}

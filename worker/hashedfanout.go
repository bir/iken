package worker

import (
	"hash/maphash"
	"sync"
)

// HashFunc converts an input to a uint value.  It must be deterministic.  See examples below.
type HashFunc[I any] func(I) uint

// HashedFanOut is a special purpose fan out that hashes the input so that the same keys are always processed
// by the same worker.  This is used to ensure guard against race conditions in non-reentrant processors.
type HashedFanOut[I any] struct {
	workerCount uint
	inputs      []chan I
	hasher      HashFunc[I]
}

func NewHashedFanOut[I any](workerCount, bufferSize uint, hasher HashFunc[I]) *HashedFanOut[I] {
	inputs := make([]chan I, workerCount)
	for i := range workerCount {
		inputs[i] = make(chan I, bufferSize)
	}

	return &HashedFanOut[I]{
		workerCount: workerCount,
		inputs:      inputs,
		hasher:      hasher,
	}
}

func (f *HashedFanOut[I]) Invoke(input I) {
	hash := f.hasher(input)
	f.inputs[hash%f.workerCount] <- input
}

func (f *HashedFanOut[I]) Close() {
	for _, c := range f.inputs {
		close(c)
	}
}

func (f *HashedFanOut[I]) Process(p ProcessorFunc[I]) {
	wg := sync.WaitGroup{}

	for _, inputChan := range f.inputs {
		wg.Add(1)

		go func(c chan I) {
			for i := range c {
				p(i)
			}

			wg.Done()
		}(inputChan)
	}

	wg.Wait()
}

type KeyFunc[I, K any] func(I) K

// StringHasher given a KeyFunc that returns a string for a given input, returns a consistent hash for the string.
func StringHasher[I any](keyFunc KeyFunc[I, string]) HashFunc[I] {
	var (
		hash maphash.Hash
		m    sync.Mutex
	)

	return func(input I) uint {
		k := keyFunc(input)

		m.Lock()
		defer m.Unlock()

		hash.Reset()
		_, _ = hash.WriteString(k)

		return uint(hash.Sum64())
	}
}

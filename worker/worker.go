package worker

import "sync"

type HashFunc[I any] func(I) int
type ProcessorFunc[I any] func(I)

type FanOut[I, O any] struct {
	workerCount int
	bufferSize  int
	inputs      chan I
}

// HashedFanOut is a special purpose fan out that hashes the input so that the same keys are always processed
// by the same worker.  This is used to ensure guard against race conditions in non-reentrant processors.
type HashedFanOut[I, O any] struct {
	FanOut[I, O]
	inputs []chan I
	hasher HashFunc[I]
}

func NewFanOut[I, O any](workerCount, bufferSize int) *FanOut[I, O] {
	f := &FanOut[I, O]{workerCount: workerCount, bufferSize: bufferSize}
	return f.init()
}

func NewHashedFanOut[I, O any](workerCount, bufferSize int, hasher HashFunc[I]) *HashedFanOut[I, O] {
	f := &HashedFanOut[I, O]{
		FanOut: FanOut[I, O]{workerCount: workerCount, bufferSize: bufferSize},
		inputs: make([]chan I, workerCount),
		hasher: hasher,
	}
	return f.init()
}

func (f *FanOut[I, O]) init() *FanOut[I, O] {
	f.inputs = make(chan I, f.bufferSize)

	return f
}

func (f *HashedFanOut[I, O]) init() *HashedFanOut[I, O] {
	for i := 0; i < f.workerCount; i++ {
		f.inputs[i] = make(chan I, f.bufferSize)
	}

	return f
}

func (f *FanOut[I, O]) Close() {
	close(f.inputs)
}

func (f *FanOut[I, O]) Invoke(input I) {
	f.inputs <- input
}

func (f *HashedFanOut[I, O]) Invoke(input I) {
	hash := f.hasher(input)
	f.inputs[hash%f.workerCount] <- input
}

func (f *HashedFanOut[I, O]) Close() {
	for _, c := range f.inputs {
		close(c)
	}
}

func (f *FanOut[I, O]) Process(p ProcessorFunc[I]) {
	for i := range f.inputs {
		p(i)
	}
}

func (f *HashedFanOut[I, O]) Process(p ProcessorFunc[I]) {
	wg := sync.WaitGroup{}

	for _, c := range f.inputs {
		wg.Add(1)
		go func(c chan I) {
			for i := range c {
				p(i)
			}

			wg.Done()
		}(c)
	}

	wg.Wait()
}

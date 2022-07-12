package worker

import "sync"

type ProcessorFunc[I any] func(I)

func NewFanOut[I any](workerCount, bufferSize uint) *FanOut[I] {
	return &FanOut[I]{workerCount: workerCount, inputs: make(chan I, bufferSize)}
}

// FanOut will fan out to `workerCount` total coroutines for processing.
type FanOut[I any] struct {
	workerCount uint
	inputs      chan I
}

// Close closes the input channels.  Invoke can not be called again after this call.
func (f *FanOut[I]) Close() {
	close(f.inputs)
}

// Invoke adds the data to the worker for processing.
func (f *FanOut[I]) Invoke(input I) {
	f.inputs <- input
}

// Process handles all inputs until the input channel is closed.
func (f *FanOut[I]) Process(p ProcessorFunc[I]) {
	wg := sync.WaitGroup{}

	for i := uint(0); i < f.workerCount; i++ {
		wg.Add(1)

		go func() {
			for i := range f.inputs {
				p(i)
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

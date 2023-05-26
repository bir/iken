package worker_test

import (
	"fmt"
	"sort"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/worker"
)

func testInts(ct int) []int {
	out := make([]int, ct)
	for i := 0; i < ct; i++ {
		out[i] = i + 1
	}

	return out
}
func TestNewFanOut(t *testing.T) {
	tests := []struct {
		name        string
		workerCount uint
		bufferSize  uint
		inputs      []int
		out         int64
	}{
		{"unbuffered single",
			1,
			0,
			testInts(5),
			15,
		},
		{"buffered x 8",
			8,
			10,
			testInts(16),
			136,
		},
		{"buffered x 2",
			2,
			10,
			testInts(16),
			136,
		},
		{"unbuffered x 16",
			16,
			0,
			testInts(16),
			136,
		},
		{"unbuffered x 2",
			2,
			0,
			testInts(16),
			136,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := worker.NewFanOut[int](tt.workerCount, tt.bufferSize)
			go func() {
				for _, i := range tt.inputs {
					w.Invoke(i)
				}
				w.Close()
			}()

			sum := int64(0)
			w.Process(func(i int) {
				atomic.AddInt64(&sum, int64(i))
			})

			assert.Equal(t, tt.out, sum)
		})
	}
}

func ExampleNewFanOut() {
	type Request struct {
		Name  string
		Index int
	}

	type Reply struct {
		Name  string
		Index int
		Size  int
	}

	inputs := []Request{{"A", 0}, {"BBBB", 1}, {"CCCCCCCCCCC", 2}}

	w := worker.NewFanOut[Request](10, 0)

	go func() {
		// Call invoke once per input data.
		for _, i := range inputs {
			w.Invoke(i)
		}

		// Call worker.FanOut.Close when all inputs are loaded.
		w.Close()
	}()

	// Unbuffered reply channel, buffer size is a tunable parameter available to the implementation
	replies := make(chan Reply)

	go func() {
		// Process and close must be executed in a separate go routine, unless the reply channel
		// is sufficiently buffered.

		w.Process(func(r Request) {
			// Do the "work".  In this example just get the size of the name.
			replies <- Reply{
				Name:  r.Name,
				Index: r.Index,
				Size:  len(r.Name),
			}
		})

		// When Process returns, all inputs have been handled.
		close(replies)
	}()

	var out []Reply
	for r := range replies {
		out = append(out, r)
	}

	// Sort the results in descending size
	sort.Slice(out, func(i int, j int) bool {
		return out[i].Size > out[j].Size
	})

	fmt.Println(out)

	// Output:
	// [{CCCCCCCCCCC 2 11} {BBBB 1 4} {A 0 1}]
}

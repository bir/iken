package worker_test

import (
	"sync/atomic"
	"testing"

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
			},
			)

			if sum != tt.out {
				t.Errorf("FanOut() = %v, want %v", sum, tt.out)
			}
		})
	}
}

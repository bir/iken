package arrays_test

import (
	"testing"

	"github.com/bir/iken/arrays"
)

func TestIndex(t *testing.T) {
	tests := []struct {
		name   string
		result int
		want   int
	}{
		{"basic int", arrays.Index(2, []int{1, 2, 3}), 1},
		{"empty int", arrays.Index(1, []int{}), -1},
		{"nil int", arrays.Index(1, nil), -1},
		{"no match int", arrays.Index(4, []int{1, 2, 3}), -1},
		{"basic string", arrays.Index("1", []string{"1", "2", "3"}), 0},
		{"basic int64", arrays.Index(1, []int64{1, 2, 3}), 0},
		{"basic ", arrays.Index(3, []int{1, 2, 3}), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result != tt.want {
				t.Errorf("got %v want %v", tt.result, tt.want)
			}
		})
	}
}

func TestIndexP(t *testing.T) {
	i := 1
	s := "1"
	tests := []struct {
		name   string
		result int
		want   int
	}{
		{"int", arrays.IndexP(&i, []int{3, 2, 1}), 2},
		{"int nil", arrays.IndexP(nil, []int{3, 2, 1}), -1},
		{"string", arrays.IndexP(&s, []string{"1", "2", "3"}), 0},
		{"string nil", arrays.IndexP(nil, []string{"1", "2", "3"}), -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result != tt.want {
				t.Errorf("got %v want %v", tt.result, tt.want)
			}
		})
	}
}

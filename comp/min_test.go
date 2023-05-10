package comp_test

import (
	"testing"

	"github.com/bir/iken/comp"
)

func TestMin(t *testing.T) {
	tests := []struct {
		name  string
		equal bool
	}{
		{"int", 1 == comp.Min(1, 99)},
		{"float32", 1 == comp.Min(1.0, 99.0)},
		{"string", "1" == comp.Min("1", "99")},
		{"uint", 1 == comp.Min(uint(1), uint(99))},
		{"uint reverse", 1 == comp.Min(uint(99), uint(1))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.equal {
				t.Error("invalid Max()")
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name  string
		equal bool
	}{
		{"int", 99 == comp.Max(1, 99)},
		{"float32", 99 == comp.Max(1.0, 99.0)},
		{"string", "99" == comp.Max("1", "99")},
		{"uint", 99 == comp.Max(uint(1), uint(99))},
		{"uint reverse", 99 == comp.Max(uint(99), uint(1))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.equal {
				t.Error("invalid Max()")
			}
		})
	}
}

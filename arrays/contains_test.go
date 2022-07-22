package arrays

import "testing"

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		result bool
		want   bool
	}{
		{"basic int", Contains(1, []int{1, 2, 3}), true},
		{"empty int", Contains(1, []int{}), false},
		{"nil int", Contains(1, nil), false},
		{"no match int", Contains(4, []int{1, 2, 3}), false},
		{"basic string", Contains("1", []string{"1", "2", "3"}), true},
		{"basic int64", Contains(1, []int64{1, 2, 3}), true},
		{"basic ", Contains(1, []int{1, 2, 3}), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result != tt.want {
				t.Errorf("got %v", tt.result)
			}
		})
	}
}

func TestContainsP(t *testing.T) {
	i := 1
	s := "1"
	tests := []struct {
		name   string
		result bool
		want   bool
	}{
		{"int", ContainsP(&i, []int{3, 2, 1}), true},
		{"string", ContainsP(&s, []string{"1", "2", "3"}), true},
		{"string nil", ContainsP(nil, []string{"1", "2", "3"}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.result != tt.want {
				t.Errorf("got %v", tt.result)
			}
		})
	}
}

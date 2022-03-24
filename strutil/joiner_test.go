package strutil

import "testing"

func mapper(s string) string {
	return s
}

func TestJoiner(t *testing.T) {
	tests := []struct {
		name    string
		keys    []string
		prefix  string
		infix   string
		postfix string
		mapper  func(string) string
		want    string
	}{
		{
			name:    "empty",
			keys:    nil,
			prefix:  "a",
			infix:   "b",
			postfix: "c",
			mapper:  mapper,
			want:    "ac",
		},
		{
			name:    "basic",
			keys:    []string{"1", "2", "3"},
			prefix:  "a",
			infix:   "b",
			postfix: "c",
			mapper:  mapper,
			want:    "a1b2b3c",
		},
		{
			name:    "single",
			keys:    []string{"1"},
			prefix:  "a",
			infix:   "b",
			postfix: "c",
			mapper:  mapper,
			want:    "a1c",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			if got := Join(test.keys, test.prefix, test.infix, test.postfix, test.mapper); got != test.want {
				t.Errorf("Join() = `%v`, want `%v`", got, test.want)
			}
		})
	}
}

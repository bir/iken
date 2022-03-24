package strutil

import "strings"

func Join[V any](ss []V, prefix, infix, postfix string, mapper func(V) string) string {
	var builder strings.Builder

	builder.WriteString(prefix)

	for i, key := range ss {
		if i > 0 {
			builder.WriteString(infix)
		}

		builder.WriteString(mapper(key))
	}

	builder.WriteString(postfix)

	return builder.String()
}

package strutil

import "strings"

func Joiner(keys []string, prefix, infix, postfix string, transform func(string) string) string {
	var builder strings.Builder

	builder.WriteString(prefix)

	for i, key := range keys {
		if i > 0 {
			builder.WriteString(infix)
		}

		builder.WriteString(transform(key))
	}

	builder.WriteString(postfix)

	return builder.String()
}

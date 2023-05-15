package httplog

import (
	"strings"
)

var RecoverBasePath = ""

func mapLine(line *string, path, prefix string) bool {
	i := strings.Index(*line, path)
	if i > 0 {
		l := prefix + (*line)[i+len(path):]
		*line = l

		return true
	}

	return false
}

func cleanPaths(line *string) {
	paths := []string{RecoverBasePath, "libexec/src/", "github.com/", "gopkg.in/"}
	prefixes := []string{"./", "\t$GO/", "\tgithub.com/", "\tgopkg.in/"}

	for n := 0; n < len(paths); n++ {
		if mapLine(line, paths[n], prefixes[n]) {
			return
		}
	}

	if (*line)[0] == '\t' {
		*line = (*line)[1:]
	}
}

func simpleLine(line, fn string) string {
	i := strings.LastIndex(fn, "/")
	if i > 0 {
		l2 := fn[:i]

		i2 := strings.LastIndex(l2, "/")
		if i2 > 0 {
			fn = fn[i2+1:]
		}
	}

	return line + " (" + fn + ")"
}

func simplifyStack(stack string, skip int) []string {
	l := strings.Split(stack, "\n")
	//	First line is goroutine ID (e.g. "goroutine 83 [running]:") - those are purged
	// The rest are pairs of lines like:
	// runtime/debug.Stack(0x0, 0x0, 0x0)
	// \t/usr/local/Cellar/go/1.11/libexec/src/runtime/debug/stack.go:24 +0xb1
	// We convert to:
	// $GO/net/http/server.go:1964 (net/http.HandlerFunc.ServeHTTP)
	result := make([]string, 0, len(l))

	var fn string

	var line string

	for i, s := range l[1+skip*2:] {
		if len(s) == 0 {
			continue
		}

		if i%2 == 0 {
			fn = s
			if strings.HasPrefix(fn, RecoverBasePath) {
				fn = strings.TrimPrefix(s, RecoverBasePath)
			}

			fn = fn[0:strings.LastIndex(fn, "(")]

			continue
		}

		line = s
		cleanPaths(&line)

		idx := strings.Index(line, " ")
		if idx > 0 {
			line = line[:idx]
		}

		r := simpleLine(line, fn)
		result = append(result, r)
	}

	return result
}

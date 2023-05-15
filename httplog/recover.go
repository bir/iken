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

func simpleLine(line, funcName string) string {
	i := strings.LastIndex(funcName, "/")
	if i > 0 {
		l2 := funcName[:i]

		i2 := strings.LastIndex(l2, "/")
		if i2 > 0 {
			funcName = funcName[i2+1:]
		}
	}

	return line + " (" + funcName + ")"
}

func simplifyStack(stack string, skip int) []string {
	lines := strings.Split(stack, "\n")
	//	First line is goroutine ID (e.g. "goroutine 83 [running]:") - those are purged
	// The rest are pairs of lines like:
	// runtime/debug.Stack(0x0, 0x0, 0x0)
	// \t/usr/local/Cellar/go/1.11/libexec/src/runtime/debug/stack.go:24 +0xb1
	// We convert to:
	// $GO/net/http/server.go:1964 (net/http.HandlerFunc.ServeHTTP)
	result := make([]string, 0, len(lines))

	var funcName string

	var line string

	for i, s := range lines[1+skip*2:] {
		if len(s) == 0 {
			continue
		}

		if i%2 == 0 {
			funcName = s
			if strings.HasPrefix(funcName, RecoverBasePath) {
				funcName = strings.TrimPrefix(s, RecoverBasePath)
			}

			funcName = funcName[0:strings.LastIndex(funcName, "(")]

			continue
		}

		line = s
		cleanPaths(&line)

		idx := strings.Index(line, " ")
		if idx > 0 {
			line = line[:idx]
		}

		r := simpleLine(line, funcName)
		result = append(result, r)
	}

	return result
}

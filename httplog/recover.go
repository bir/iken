package httplog

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/rs/zerolog"

	"github.com/bir/iken/httputil"
)

// ErrInternal is the default error returned from a panic.
var ErrInternal = errors.New("internal error")

// RecoverLogger returns a handler that call initializes Op in the context, and logs each request.
func RecoverLogger(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := log.With().Logger().WithContext(r.Context())

			defer func() {
				rErr := recover()
				if rErr != nil {
					var err error
					switch t := rErr.(type) {
					case string:
						err = fmt.Errorf("%v: %w", t, ErrInternal)
					case error:
						err = t
					default:
						err = ErrInternal
					}
					s := string(debug.Stack())

					zerolog.Ctx(ctx).Err(err).Strs(Stack, simplifyStack(s, stackSkip)).Msg("Panic")

					httputil.HTTPInternalServerError(w, r)
				}
			}()

			if next != nil {
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}

var RecoverBasePath = initBasePath()

func initBasePath() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	return buildInfo.Main.Path
}

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
	paths := []string{RecoverBasePath, "libexec/src/", "github.com/", "gopkg.in/", "x64/src"}
	prefixes := []string{"./", "\t$GO/", "\tgithub.com/", "\tgopkg.in/", "\t$GO/"}

	for n := 0; n < len(paths); n++ {
		if mapLine(line, paths[n], prefixes[n]) {
			return
		}
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
			funcName = s[0:strings.LastIndex(s, "(")]

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

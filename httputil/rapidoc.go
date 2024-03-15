package httputil

import (
	"bytes"
	"html/template"
	"net/http"
)

// RapiDocOpts configures the RapiDoc middlewares.
type RapiDocOpts struct {
	// BasePath for the UI path, defaults to: /
	BasePath string
	// Path combines with BasePath for the full UI path, defaults to: docs
	Path string
	// SpecURL the url to find the spec for
	SpecURL string
	// RapiDocURL for the js that generates the rapidoc site, defaults to: https://unpkg.com/rapidoc/dist/rapidoc-min.js
	RapiDocURL string
	// Title for the documentation site, default to: API documentation
	Title string
}

// Defaults for all options.
func (r *RapiDocOpts) Defaults() {
	if r.BasePath == "" {
		r.BasePath = "/"
	}

	if r.Path == "" {
		r.Path = "docs"
	}

	if r.SpecURL == "" {
		r.SpecURL = "/swagger.json"
	}

	if r.RapiDocURL == "" {
		r.RapiDocURL = rapidocLatest
	}

	if r.Title == "" {
		r.Title = "API Documentation"
	}
}

// RapiDoc creates a handler to serve a documentation site for a swagger spec.
// This allows for altering the spec before starting the http listener.
func RapiDoc(opts RapiDocOpts) func(w http.ResponseWriter, r *http.Request) {
	opts.Defaults()

	tmpl := template.Must(template.New("rapidoc").Parse(rapidocTemplate))

	buf := bytes.NewBuffer(nil)
	_ = tmpl.Execute(buf, opts)
	b := buf.Bytes()

	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		_, _ = w.Write(b)
	}
}

const (
	rapidocLatest   = "https://unpkg.com/rapidoc/dist/rapidoc-min.js"
	rapidocTemplate = `<!doctype html>
<html>
<head>
  <title>{{ .Title }}</title>
  <meta charset="utf-8">
  <script type="module" src="{{ .RapiDocURL }}"></script>
</head>
<body>
  <rapi-doc spec-url="{{ .SpecURL }}" render-style = "read"
  show-header = 'false'
  allow-server-selection = 'false'
  theme = "dark"></rapi-doc>
</body>
</html>
`
)

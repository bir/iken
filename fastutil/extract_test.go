package fastutil_test

import (
	"strings"
	"testing"

	"github.com/bir/iken/fastutil"
	"github.com/valyala/fasthttp"
)

func makeCtxWithQuery(q string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("GET ?" + q)

	return ctx
}

func makeCtxWithPathParam(k, v string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.SetUserValue(k, v)

	return ctx
}

func TestQueryString(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		key     string
		want    string
		wantErr bool
	}{
		{"basic", "test=foo", "test", "foo", false},
		{"err", "test=foo", "foo", "", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := makeCtxWithQuery(test.query)
			got, err := fastutil.QueryString(ctx, test.key)
			if (err != nil) != test.wantErr {
				t.Errorf("error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("got = %v, want %v", got, test.want)
			}
		})
	}
}

func TestQueryStringOptional(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		key     string
		want    string
		wantNil bool
	}{
		{"basic", "test=foo", "test", "foo", false},
		{"nil", "test=foo", "foo", "", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := makeCtxWithQuery(test.query)
			got := fastutil.QueryStringOptional(ctx, test.key)
			if got == nil {
				if !test.wantNil {
					t.Errorf("error = %v, wantNil %v", got, test.wantNil)
				}
				return
			}

			if *got != test.want {
				t.Errorf("got = %v, want %v", got, test.want)
			}
		})
	}
}

func TestQueryStringsOptional(t *testing.T) {
	tests := []struct {
		name  string
		query string
		key   string
		want  string
	}{
		{"single", "test=foo", "test", "foo"},
		{"nil", "test=foo", "foo", ""},
		{"two", "test=foo&test=foo2", "test", "foo,foo2"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := makeCtxWithQuery(test.query)
			got := fastutil.QueryStringsOptional(ctx, test.key)

			if strings.Join(got, ",") != test.want {
				t.Errorf("got = %v, want %v", got, test.want)
			}
		})
	}
}

func TestQueryInt32Optional(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		key     string
		want    int32
		wantNil bool
	}{
		{"151", "test=151&x=2", "test", 151, false},
		{"nil", "test=foo", "foo", 0, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := makeCtxWithQuery(test.query)

			got := fastutil.QueryInt32Optional(ctx, test.key)
			if got == nil {
				if !test.wantNil {
					t.Errorf("error = %v, wantNil %v", got, test.wantNil)
				}
				return
			}

			if *got != test.want {
				t.Errorf("got = %v, want %v", got, test.want)
			}
		})
	}
}

func TestQueryInt64Optional(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		key     string
		want    int64
		wantNil bool
	}{
		{"151", "test=151&x=2", "test", 151, false},
		{"nil", "test=foo", "foo", 0, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := makeCtxWithQuery(test.query)

			got := fastutil.QueryInt64Optional(ctx, test.key)
			if got == nil {
				if !test.wantNil {
					t.Errorf("error = %v, wantNil %v", got, test.wantNil)
				}
				return
			}

			if *got != test.want {
				t.Errorf("got = %v, want %v", got, test.want)
			}
		})
	}
}

func TestQueryBool(t *testing.T) {
	tests := []struct {
		name  string
		query string
		key   string
		want  bool
	}{
		{"true", "test=true&x=2", "test", true},
		{"false", "test=false", "test", false},
		{"default", "test=foo", "foo", false},
		{"1", "test=1", "test", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := makeCtxWithQuery(test.query)

			got, _ := fastutil.QueryBool(ctx, test.key)
			if got != test.want {
				t.Errorf("got = %v, want %v", got, test.want)
			}
		})
	}
}

func TestPathInt64(t *testing.T) {
	tests := []struct {
		name    string
		inKey   string
		param   string
		key     string
		want    int64
		wantErr bool
	}{
		{"2", "x", "2", "x", 2, false},
		{"empty", "x", "", "x", 0, true},
		{"missing", "y", "2", "x", 2, true},
		{"bad string", "x", "a2", "x", 2, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := makeCtxWithPathParam(test.inKey, test.param)

			got, err := fastutil.PathInt64(ctx, test.key)
			if err != nil {
				if !test.wantErr {
					t.Errorf("got = %v, want %v", err, test.want)
				}
				return
			}

			if got != test.want {
				t.Errorf("got = %v, want %v", got, test.want)
			}
		})
	}
}

func TestErrRequired(t *testing.T){
	if fastutil.ErrRequired.Error() != "required" {
		t.Errorf("got = %v, want %v", fastutil.ErrRequired.Error(),"required")
	}
}
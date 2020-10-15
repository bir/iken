package fastctx_test

import (
	"encoding/json"
	"testing"

	"github.com/bir/iken/fastctx"
	"github.com/valyala/fasthttp"
)

func TestDebugMap(t *testing.T) {
	c := fasthttp.RequestCtx{}

	got := fastctx.DebugMap(&c)

	j, err := json.Marshal(got)
	if err != nil {
		t.Error(err)
		return
	}
	want := `{"connID":0,"ip":"0.0.0.0:0","requestID":0,"requestNum":0}`
	if string(j) != want {
		t.Errorf("want `%#v`, got `%v`", want, string(j))
	}

	got = fastctx.DebugMap(nil)
	if got != nil {
		t.Errorf("want `%#v`, got `%v`", nil, got)
	}

}

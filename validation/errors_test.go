package validation_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/validation"
)

func TestErrors_Add(t *testing.T) {
	tests := []struct {
		name     string
		ee       validation.Errors
		field    string
		msg      string
		want     string
		wantJson string
	}{
		{"empty", nil, "", "", "", `{}`},
		{"empty add", nil, "test", "bad", "test: bad.", `{"test":["bad"]}`},
		{"add new", *(&validation.Errors{}).Add("a", errors.New("b")), "test", "bad", "a: b; test: bad.", `{"a":["b"],"test":["bad"]}`},
		{"add new string", *(&validation.Errors{}).Add("a", "b"), "test", "bad", "a: b; test: bad.", `{"a":["b"],"test":["bad"]}`},
		{"add existing", *(&validation.Errors{}).Add("a", errors.New("b")), "a", "x", "a: b, x.", `{"a":["b","x"]}`},
		{"add User",
			*(&validation.Errors{}).Add("a",
				validation.Error{Message: "PUBLIC", Source: errors.New("PRIVATE")}),
			"a",
			"x",
			"a: PUBLIC: PRIVATE, x.",
			`{"a":["PUBLIC","x"]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &tt.ee
			if tt.field != "" {
				got = tt.ee.Add(tt.field, errors.New(tt.msg))
			}
			assert.Equal(t, tt.want, got.Error())

			b, err := json.Marshal(got.Fields())
			assert.Nil(t, err)

			assert.Equal(t, tt.wantJson, string(b))
		})
	}
}

func TestErrors_Error(t *testing.T) {
	var ee validation.Errors
	assert.Empty(t, ee.Error())
}

func TestErrors_GetErr(t *testing.T) {
	var ee validation.Errors
	assert.Nil(t, ee.GetErr())

	_ = ee.Add("a", errors.New("b"))

	assert.NotNil(t, ee.GetErr())
	assert.Equal(t, "a: b.", ee.GetErr().Error())
}

func TestErrors_New(t *testing.T) {
	err := validation.New("a", "b")
	assert.NotEmpty(t, err)
	assert.Equal(t, "a: b.", err.Error())
}

func TestErrors_NewError(t *testing.T) {
	err := validation.NewError("a", errors.New("b"))
	assert.NotEmpty(t, err)
	assert.Equal(t, "a: b.", err.Error())
}

func TestJoin(t *testing.T) {
	s := validation.Join(nil, "|")
	assert.Empty(t, s)
}

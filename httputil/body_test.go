package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bir/iken/validation"
)

type TestObject struct {
	ID    string
	Count int
}

func (p *TestObject) UnmarshalJSON(b []byte) error {
	var requiredCheck map[string]any

	if err := json.Unmarshal(b, &requiredCheck); err != nil {
		return validation.Error{err.Error(), fmt.Errorf("TestObject.UnmarshalJSON Required: `%v`: %w", string(b), err)}
	}

	var validationErrors validation.Errors

	if _, ok := requiredCheck["ID"]; !ok {
		return validationErrors.Add("message_id", "missing required field")
	}

	type TestObjectJSON TestObject
	var parseObject TestObjectJSON

	if err := json.Unmarshal(b, &parseObject); err != nil {
		return validation.Error{err.Error(), fmt.Errorf("Message.UnmarshalJSON: `%v`: %w", string(b), err)}
	}

	*p = TestObject(parseObject)

	return nil
}

func strP(s string) *string {
	return &s
}

type BadIOReader struct {
	err error
}

func (r *BadIOReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

func TestGetJSONBody(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		body    any
		want    any
		wantErr bool
	}{
		{"no body", nil, nil, nil, true},
		{"string", bytes.NewBufferString(`"foo"`), strP(""), strP("foo"), false},
		{"invalid json", bytes.NewBufferString(`{"foo"`), strP(""), strP(""), true},
		{"empty", bytes.NewBufferString(``), strP(""), strP(""), true},
		{"EOF", &BadIOReader{io.EOF}, strP(""), strP(""), true},
		{"read error", &BadIOReader{io.ErrClosedPipe}, strP(""), strP(""), true},
		{"null body", bytes.NewBufferString(`null`), &TestObject{}, &TestObject{}, true},
		{"validation error - bad ID type", bytes.NewBufferString(`{"ID":1}`), &TestObject{}, &TestObject{}, true},
		{"validations error - no ID", bytes.NewBufferString(`{}`), &TestObject{}, &TestObject{}, true},
		{"good", bytes.NewBufferString(`{"ID":"1"}`), &TestObject{}, &TestObject{ID: "1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GetJSONBody(tt.r, tt.body)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, tt.body)
		})
	}
}

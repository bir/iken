package params

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetInt32(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		required bool
		want     int32
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=123", nil), "foo", true, 123, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, 0, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, 0, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, 0, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetInt32(tt.r, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		required bool
		want     int
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=123", nil), "foo", true, 123, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, 0, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, 0, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, 0, true, false},
		{"max", httptest.NewRequest("GET", "/BAR?foo=9223372036854775807", nil), "foo", true, 9223372036854775807, false, true},
		{"over max", httptest.NewRequest("GET", "/BAR?foo=19223372036854775807", nil), "foo", true, 0, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetInt(tt.r, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetInt64(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		required bool
		want     int64
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=123", nil), "foo", true, 123, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, 0, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, 0, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, 0, true, false},
		{"max", httptest.NewRequest("GET", "/BAR?foo=9223372036854775807", nil), "foo", true, 9223372036854775807, false, true},
		{"over max", httptest.NewRequest("GET", "/BAR?foo=19223372036854775807", nil), "foo", true, 0, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetInt64(tt.r, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		required bool
		want     bool
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=true", nil), "foo", true, true, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, false, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, false, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, false, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetBool(tt.r, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetInt32Array(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		required bool
		want     []int32
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=123", nil), "foo", true, []int32{123}, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, nil, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, nil, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, nil, true, false},
		{"large", httptest.NewRequest("GET", "/BAR?foo=1,2,3,4", nil), "foo", true, []int32{1, 2, 3, 4}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetInt32Array(tt.r, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetTime(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		required bool
		want     time.Time
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=2006-01-02T15:04:05Z", nil), "foo", true, time.Date(2006, 01, 02, 15, 4, 5, 0, time.UTC), false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, time.Time{}, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, time.Time{}, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=200601021504050700", nil), "foo", true, time.Time{}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetTime(tt.r, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetUUID(t *testing.T) {
	testUUID, _ := uuid.Parse("48ab873f-d4fc-4e2b-bf92-9440e431ff54")

	tests := []struct {
		name     string
		r        *http.Request
		param    string
		required bool
		want     uuid.UUID
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=48ab873f-d4fc-4e2b-bf92-9440e431ff54", nil), "foo", true, testUUID, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, uuid.UUID{}, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, uuid.UUID{}, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, uuid.UUID{}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetUUID(tt.r, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestURLParam(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	ctx := chi.NewRouteContext()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
	ctx.URLParams.Add("id", "12345")

	got, ok, err := GetInt(r, "id", true)

	assert.Nil(t, err)
	assert.NotEmpty(t, got)
	assert.True(t, ok)
	assert.Equal(t, got, 12345)
}

type TestEnum int8

const (
	testEnumUnknown TestEnum = iota
	testEnumA
	testEnumB
	testEnumC
)

func NewTestEnum(name string) TestEnum {
	switch name {
	case "aaa":
		return testEnumA
	case "bbb":
		return testEnumB
	case "ccc":
		return testEnumC
	}

	return TestEnum(0)
}

func TestGetEnum(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		required bool
		want     TestEnum
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=bbb", nil), "foo", true, testEnumB, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, testEnumUnknown, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, testEnumUnknown, false, false},
		{"bad value", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, testEnumUnknown, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetEnum(tt.r, tt.param, tt.required, NewTestEnum)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetEnumArray(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		required bool
		want     []TestEnum
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=bbb", nil), "foo", true, []TestEnum{testEnumB}, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, nil, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, nil, false, false},
		{"bad value", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, []TestEnum{testEnumUnknown}, false, true},
		{"all", httptest.NewRequest("GET", "/BAR?foo=aaa,bbb,ccc", nil), "foo", true, []TestEnum{testEnumA, testEnumB, testEnumC}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetEnumArray(tt.r, tt.param, tt.required, NewTestEnum)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

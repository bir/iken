package params

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	newHeaderRequest := func(key, value string) *http.Request {
		r := httptest.NewRequest("GET", "/ping", nil)
		if key != "" {
			r.Header.Set(key, value)
		}

		return r
	}

	tests := []struct {
		name     string
		r        *http.Request
		param    string
		required bool
		want     string
		wantErr  bool
		wantOk   bool
	}{
		{" header required present", newHeaderRequest("foo", "123"), "foo", true, "123", false, true},
		{" header required missing", newHeaderRequest("", ""), "foo", true, "", true, false},
		{" header not required present", newHeaderRequest("foo", "123"), "foo", false, "123", false, true},
		{" header not required missing", newHeaderRequest("", ""), "foo", false, "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetString(tt.r, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

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
		{"simple", httptest.NewRequest("GET", "/BAR?foo=2006-01-02T15:04:05Z", nil), "foo", true, time.Date(2006, 0o1, 0o2, 15, 4, 5, 0, time.UTC), false, true},
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
	r.SetPathValue("id", "12345")

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
func TestGetStringFrom(t *testing.T) {
	newHeaderRequest := func(key, value string) *http.Request {
		r := httptest.NewRequest("GET", "/ping", nil)
		if key != "" {
			r.Header.Set(key, value)
		}

		return r
	}

	newMultiSourceRequest := func(key string) *http.Request {
		r := httptest.NewRequest("GET", "/ping?"+key+"=queryval", nil)
		r.Header.Set(key, "headerval")
		r.AddCookie(&http.Cookie{Name: key, Value: "cookieval"})
		r.SetPathValue(key, "pathval")
		return r
	}

	tests := []struct {
		name     string
		r        *http.Request
		param    string
		source   ParamSource
		required bool
		want     string
		wantErr  bool
		wantOk   bool
	}{
		{" header required present", newHeaderRequest("foo", "123"), "foo", ParamHeader, true, "123", false, true},
		{" header required missing", newHeaderRequest("", ""), "foo", ParamHeader, true, "", true, false},
		{" header not required present", newHeaderRequest("foo", "123"), "foo", ParamHeader, false, "123", false, true},
		{" header not required missing", newHeaderRequest("", ""), "foo", ParamHeader, false, "", false, false},
		{" fetch header specifically", newMultiSourceRequest("foo"), "foo", ParamHeader, false, "headerval", false, true},
		{" fetch query specifically", newMultiSourceRequest("foo"), "foo", ParamQuery, false, "queryval", false, true},
		{" fetch cookie specifically", newMultiSourceRequest("foo"), "foo", ParamCookie, false, "cookieval", false, true},
		{" fetch cookie specifically", newMultiSourceRequest("foo"), "foo", ParamPath, false, "pathval", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetStringFrom(tt.r, tt.param, tt.source, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetInt32From(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		source   ParamSource
		required bool
		want     int32
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=123", nil), "foo", ParamQuery, true, 123, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", ParamQuery, true, 0, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", ParamQuery, false, 0, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", ParamQuery, true, 0, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetInt32From(tt.r, tt.param, tt.source, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetIntFrom(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		source   ParamSource
		required bool
		want     int
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=123", nil), "foo", ParamQuery, true, 123, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", ParamQuery, true, 0, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", ParamQuery, false, 0, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", ParamQuery, true, 0, true, false},
		{"max", httptest.NewRequest("GET", "/BAR?foo=9223372036854775807", nil), "foo", ParamQuery, true, 9223372036854775807, false, true},
		{"over max", httptest.NewRequest("GET", "/BAR?foo=19223372036854775807", nil), "foo", ParamQuery, true, 0, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetIntFrom(tt.r, tt.param, tt.source, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetInt64From(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		source   ParamSource
		required bool
		want     int64
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=123", nil), "foo", ParamQuery, true, 123, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", ParamQuery, true, 0, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", ParamQuery, false, 0, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", ParamQuery, true, 0, true, false},
		{"max", httptest.NewRequest("GET", "/BAR?foo=9223372036854775807", nil), "foo", ParamQuery, true, 9223372036854775807, false, true},
		{"over max", httptest.NewRequest("GET", "/BAR?foo=19223372036854775807", nil), "foo", ParamQuery, true, 0, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetInt64From(tt.r, tt.param, tt.source, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetBoolFrom(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		source   ParamSource
		required bool
		want     bool
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=true", nil), "foo", ParamQuery, true, true, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", ParamQuery, true, false, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", ParamQuery, false, false, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", ParamQuery, true, false, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetBoolFrom(tt.r, tt.param, tt.source, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetInt32ArrayFrom(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		source   ParamSource
		required bool
		want     []int32
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=123", nil), "foo", ParamQuery, true, []int32{123}, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", ParamQuery, true, nil, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", ParamQuery, false, nil, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", ParamQuery, true, nil, true, false},
		{"large", httptest.NewRequest("GET", "/BAR?foo=1,2,3,4", nil), "foo", ParamQuery, true, []int32{1, 2, 3, 4}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetInt32ArrayFrom(tt.r, tt.param, tt.source, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetTimeFrom(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		source   ParamSource
		required bool
		want     time.Time
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=2006-01-02T15:04:05Z", nil), "foo", ParamQuery, true, time.Date(2006, 0o1, 0o2, 15, 4, 5, 0, time.UTC), false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", ParamQuery, true, time.Time{}, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", ParamQuery, false, time.Time{}, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=200601021504050700", nil), "foo", ParamQuery, true, time.Time{}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetTimeFrom(tt.r, tt.param, tt.source, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetUUIDFrom(t *testing.T) {
	testUUID, _ := uuid.Parse("48ab873f-d4fc-4e2b-bf92-9440e431ff54")

	tests := []struct {
		name     string
		r        *http.Request
		param    string
		source   ParamSource
		required bool
		want     uuid.UUID
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=48ab873f-d4fc-4e2b-bf92-9440e431ff54", nil), "foo", ParamQuery, true, testUUID, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", ParamQuery, true, uuid.UUID{}, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", ParamQuery, false, uuid.UUID{}, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", ParamQuery, true, uuid.UUID{}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetUUIDFrom(tt.r, tt.param, tt.source, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetEnumFrom(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		source   ParamSource
		required bool
		want     TestEnum
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=bbb", nil), "foo", ParamQuery, true, testEnumB, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", ParamQuery, true, testEnumUnknown, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", ParamQuery, false, testEnumUnknown, false, false},
		{"bad value", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", ParamQuery, true, testEnumUnknown, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetEnumFrom(tt.r, tt.param, tt.source, tt.required, NewTestEnum)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

func TestGetEnumArrayFrom(t *testing.T) {
	tests := []struct {
		name     string
		r        *http.Request
		param    string
		source   ParamSource
		required bool
		want     []TestEnum
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=bbb", nil), "foo", ParamQuery, true, []TestEnum{testEnumB}, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", ParamQuery, true, nil, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", ParamQuery, false, nil, false, false},
		{"bad value", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", ParamQuery, true, []TestEnum{testEnumUnknown}, false, true},
		{"all", httptest.NewRequest("GET", "/BAR?foo=aaa,bbb,ccc", nil), "foo", ParamQuery, true, []TestEnum{testEnumA, testEnumB, testEnumC}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetEnumArrayFrom(tt.r, tt.param, tt.source, tt.required, NewTestEnum)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
			assert.Equal(t, tt.wantOk, ok, "ok")
		})
	}
}

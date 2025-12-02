package params

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func NewQueryRequest(method, target, key, value string) *http.Request {
	r := httptest.NewRequest(method, target+"?"+key+"="+value, nil)
	return r
}

func NewPathRequest(method, target, key, value string) *http.Request {
	r := httptest.NewRequest(method, target, nil)
	r.SetPathValue(key, value)
	return r
}

func NewHeaderRequest(method, target, key, value string) *http.Request {
	r := httptest.NewRequest(method, target, nil)
	r.Header.Add(key, value)
	return r
}

func NewCookieRequest(method, target, key, value string) *http.Request {
	r := httptest.NewRequest(method, target, nil)
	r.AddCookie(&http.Cookie{Name: key, Value: value})
	return r
}

func NewMultiSourceRequest(method, target, key string, values [4]string) *http.Request {
	r := httptest.NewRequest(method, target+"?"+key+"="+values[ParamSourceQuery], nil)
	r.Header.Set(key, values[ParamSourceHeader])
	r.AddCookie(&http.Cookie{Name: key, Value: values[ParamSourceCookie]})
	r.SetPathValue(key, values[ParamSourcePath])
	return r
}

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

func TestGetUUIDArray(t *testing.T) {
	id1, id2, id3 := uuid.New(), uuid.New(), uuid.New()

	tests := []struct {
		name     string
		r        *http.Request
		param    string
		required bool
		want     []uuid.UUID
		wantErr  bool
		wantOk   bool
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo="+id1.String(), nil), "foo", true, []uuid.UUID{id1}, false, true},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, nil, true, false},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, nil, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, nil, true, false},
		{"large", httptest.NewRequest("GET", fmt.Sprintf("/BAR?foo=%s,%s,%s", id1.String(), id2.String(), id3.String()), nil), "foo", true, []uuid.UUID{id1, id2, id3}, false, true},
		{"large repeated", httptest.NewRequest("GET", fmt.Sprintf("/BAR?foo=%s&foo=%s,,,%s", id1.String(), id2.String(), id3.String()), nil), "foo", true, []uuid.UUID{id1, id2, id3}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok, err := GetUUIDArray(tt.r, tt.param, tt.required)

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

type ParamSource int

const (
	ParamSourcePath ParamSource = iota
	ParamSourceQuery
	ParamSourceHeader
	ParamSourceCookie
)

var (
	ParamSources     = []ParamSource{ParamSourcePath, ParamSourceQuery, ParamSourceHeader, ParamSourceCookie}
	ParamSourceNames = []string{"Path", "Query", "Header", "Cookie"}
)

// This abomination of a test function exists to run the same tests against all the different types
// that can be retrieved from a param and all the different ways a param can be passed. The
// ugliness in typeList and everything downstream of that is necessary to get the polymorpism
// required to run against different types.
func TestMatrix(t *testing.T) {
	testUUID, _ := uuid.Parse("48ab873f-d4fc-4e2b-bf92-9440e431ff54")
	RequestFunctions := map[ParamSource]func(method string, target string, key string, value string) *http.Request{
		ParamSourcePath:   NewPathRequest,
		ParamSourceQuery:  NewQueryRequest,
		ParamSourceHeader: NewHeaderRequest,
		ParamSourceCookie: NewCookieRequest,
	}

	typeList := []struct {
		Name              string
		TestValue         any
		TestValueAsString string
		// Go does not allow you to pass a func() something as a func() any, so wrappers will be needed :(
		Methods map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error)
	}{
		{
			Name:              "String",
			TestValue:         "foo,bar",
			TestValueAsString: "foo,bar",
			Methods: map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error){
				ParamSourcePath: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetStringPath(r, name, required)
				},
				ParamSourceQuery: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetStringQuery(r, name, required)
				},
				ParamSourceHeader: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetStringHeader(r, name, required)
				},
				ParamSourceCookie: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetStringCookie(r, name, required)
				},
			},
		},
		{
			Name:              "Int",
			TestValue:         123,
			TestValueAsString: "123",
			Methods: map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error){
				ParamSourcePath: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetIntPath(r, name, required)
				},
				ParamSourceQuery: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetIntQuery(r, name, required)
				},
				ParamSourceHeader: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetIntHeader(r, name, required)
				},
				ParamSourceCookie: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetIntCookie(r, name, required)
				},
			},
		},
		{
			Name:              "Int32",
			TestValue:         int32(123),
			TestValueAsString: "123",
			Methods: map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error){
				ParamSourcePath: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetInt32Path(r, name, required)
				},
				ParamSourceQuery: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetInt32Query(r, name, required)
				},
				ParamSourceHeader: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetInt32Header(r, name, required)
				},
				ParamSourceCookie: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetInt32Cookie(r, name, required)
				},
			},
		},
		{
			Name:              "Int64",
			TestValue:         int64(123),
			TestValueAsString: "123",
			Methods: map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error){
				ParamSourcePath: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetInt64Path(r, name, required)
				},
				ParamSourceQuery: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetInt64Query(r, name, required)
				},
				ParamSourceHeader: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetInt64Header(r, name, required)
				},
				ParamSourceCookie: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetInt64Cookie(r, name, required)
				},
			},
		},
		{
			Name:              "Bool",
			TestValue:         true,
			TestValueAsString: "true",
			Methods: map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error){
				ParamSourcePath: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetBoolPath(r, name, required)
				},
				ParamSourceQuery: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetBoolQuery(r, name, required)
				},
				ParamSourceHeader: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetBoolHeader(r, name, required)
				},
				ParamSourceCookie: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetBoolCookie(r, name, required)
				},
			},
		},
		{
			Name:              "Time",
			TestValue:         time.Date(2006, 0o1, 0o2, 15, 4, 5, 0, time.UTC),
			TestValueAsString: "2006-01-02T15:04:05Z",
			Methods: map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error){
				ParamSourcePath: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetTimePath(r, name, required)
				},
				ParamSourceQuery: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetTimeQuery(r, name, required)
				},
				ParamSourceHeader: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetTimeHeader(r, name, required)
				},
				ParamSourceCookie: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetTimeCookie(r, name, required)
				},
			},
		},
		{
			Name:              "UUID",
			TestValue:         testUUID,
			TestValueAsString: "48ab873f-d4fc-4e2b-bf92-9440e431ff54",
			Methods: map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error){
				ParamSourcePath: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetUUIDPath(r, name, required)
				},
				ParamSourceQuery: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetUUIDQuery(r, name, required)
				},
				ParamSourceHeader: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetUUIDHeader(r, name, required)
				},
				ParamSourceCookie: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetUUIDCookie(r, name, required)
				},
			},
		},
		{
			Name:              "Enum",
			TestValue:         testEnumA,
			TestValueAsString: "aaa",
			Methods: map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error){
				ParamSourcePath: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetEnumPath(r, name, required, NewTestEnum)
				},
				ParamSourceQuery: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetEnumQuery(r, name, required, NewTestEnum)
				},
				ParamSourceHeader: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetEnumHeader(r, name, required, NewTestEnum)
				},
				ParamSourceCookie: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetEnumCookie(r, name, required, NewTestEnum)
				},
			},
		},
		{
			Name:              "StringArray",
			TestValue:         []string{"a", "b"},
			TestValueAsString: "a,b",
			Methods: map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error){
				ParamSourcePath: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetStringArrayPath(r, name, required)
				},
				ParamSourceQuery: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetStringArrayQuery(r, name, required)
				},
				ParamSourceHeader: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetStringArrayHeader(r, name, required)
				},
				ParamSourceCookie: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetStringArrayCookie(r, name, required)
				},
			},
		},
		{
			Name:              "Int32Array",
			TestValue:         []int32{123, 456},
			TestValueAsString: "123,456",
			Methods: map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error){
				ParamSourcePath: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetInt32ArrayPath(r, name, required)
				},
				ParamSourceQuery: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetInt32ArrayQuery(r, name, required)
				},
				ParamSourceHeader: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetInt32ArrayHeader(r, name, required)
				},
				ParamSourceCookie: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetInt32ArrayCookie(r, name, required)
				},
			},
		},
		{
			Name:              "UUIDArray",
			TestValue:         []uuid.UUID{uuid.MustParse("902da57e-3e3a-470a-b821-0cd140a7f442"), uuid.MustParse("bcdcf46c-1baf-4e95-b607-97cf4aca1877")},
			TestValueAsString: "902da57e-3e3a-470a-b821-0cd140a7f442,bcdcf46c-1baf-4e95-b607-97cf4aca1877",
			Methods: map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error){
				ParamSourcePath: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetUUIDArrayPath(r, name, required)
				},
				ParamSourceQuery: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetUUIDArrayQuery(r, name, required)
				},
				ParamSourceHeader: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetUUIDArrayHeader(r, name, required)
				},
				ParamSourceCookie: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetUUIDArrayCookie(r, name, required)
				},
			},
		},
		{
			Name:              "EnumArray",
			TestValue:         []TestEnum{testEnumA, testEnumB},
			TestValueAsString: "aaa,bbb",
			Methods: map[ParamSource]func(r *http.Request, name string, required bool) (any, bool, error){
				ParamSourcePath: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetEnumArrayPath(r, name, required, NewTestEnum)
				},
				ParamSourceQuery: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetEnumArrayQuery(r, name, required, NewTestEnum)
				},
				ParamSourceHeader: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetEnumArrayHeader(r, name, required, NewTestEnum)
				},
				ParamSourceCookie: func(r *http.Request, name string, required bool) (any, bool, error) {
					return GetEnumArrayCookie(r, name, required, NewTestEnum)
				},
			},
		},
	}

	for _, typeInfo := range typeList {
		for _, source := range ParamSources {
			multiSourceValues := [4]string{"a", "b", "c", "d"}
			multiSourceValues[source] = typeInfo.TestValueAsString

			perTypeTests := []struct {
				Name     string
				Request  *http.Request
				Required bool
				Want     any
				WantErr  bool
				WantOk   bool
				Method   func(r *http.Request, name string, required bool) (any, bool, error)
			}{
				{
					Name:     "required present",
					Request:  RequestFunctions[source]("GET", "/BAR", "foo", typeInfo.TestValueAsString),
					Method:   typeInfo.Methods[source],
					Required: true,
					Want:     typeInfo.TestValue,
					WantOk:   true,
				},
				{
					Name:     "required missing",
					Request:  RequestFunctions[source]("GET", "/BAR", "", ""),
					Method:   typeInfo.Methods[source],
					Required: true,
					WantErr:  true,
				},
				{
					Name:     "optional missing",
					Request:  RequestFunctions[source]("GET", "/BAR", "", ""),
					Method:   typeInfo.Methods[source],
					Required: false,
					WantOk:   false,
				},
				{
					Name:    "ignore other sources",
					Request: NewMultiSourceRequest("GET", "/BAR", "foo", multiSourceValues),
					Method:  typeInfo.Methods[source],
					Want:    typeInfo.TestValue,
					WantOk:  true,
				},
			}

			for _, tt := range perTypeTests {
				t.Run("Get"+typeInfo.Name+ParamSourceNames[source]+"_"+tt.Name, func(t *testing.T) {
					got, ok, err := tt.Method(tt.Request, "foo", tt.Required)

					if tt.WantErr {
						assert.Error(t, err)
					} else if tt.WantOk {
						assert.NoError(t, err)
						assert.True(t, ok)
						assert.Equal(t, tt.Want, got)
					} else {
						assert.NoError(t, err)
						assert.False(t, ok)
					}
				})
			}
		}
	}
}

package forms

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetFile(t *testing.T) {
	const expectedKey = "foo"
	const expectedContent = "123"
	const expectedFilename = "test.txt"

	newRequest := func(notMultipart bool) (*http.Request, error) {
		if notMultipart {
			return http.NewRequest("GET", "/ping", nil)
		}

		var data bytes.Buffer
		w := multipart.NewWriter(&data)

		fw, err := w.CreateFormFile(expectedKey, expectedFilename)
		if err != nil {
			return nil, fmt.Errorf("error creating field: %w", err)
		}

		_, err = io.Copy(fw, strings.NewReader(expectedContent))
		if err != nil {
			return nil, fmt.Errorf("error copying value for field: %w", err)
		}

		err = w.Close()
		if err != nil {
			return nil, fmt.Errorf("error closing writer: for field: %w", err)
		}

		req := httptest.NewRequest("POST", "/ping", &data)
		req.Header.Set("Content-Type", w.FormDataContentType())

		return req, nil
	}

	tests := []struct {
		name         string
		key          string
		notMultipart bool
		required     bool
		want         string
		errMsg       string
	}{
		{name: " required present", key: "foo", required: true, want: "123"},
		{name: " required missing", key: "foo2", required: true, errMsg: ErrNotFound.Error()},
		{name: " not required present", key: "foo", want: "123"},
		{name: " not required missing", key: "foo2"},
		{name: " not multipart", key: "foo", notMultipart: true, errMsg: "ToFormFile: request Content-Type isn't multipart/form-data"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := newRequest(tt.notMultipart)
			assert.NoError(t, err)

			f, ok, err := GetFile(r, tt.key, tt.required)
			if tt.errMsg != "" {
				assert.ErrorContains(t, err, tt.errMsg)
				return
			}

			assert.NoError(t, err)

			if ok {
				defer f.File.Close()
			}

			if tt.want == "" {
				assert.Nil(t, f)
				return
			}

			assert.Equal(t, expectedFilename, f.Filename)
			assert.Equal(t, int64(len(expectedContent)), f.Size)

			buf := new(strings.Builder)
			_, err = io.Copy(buf, f.File)
			assert.NoError(t, err)
			assert.Equal(t, len(expectedContent), len(buf.String()))
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestGetString(t *testing.T) {
	newFormRequest := func(key, value string) *http.Request {
		var data string
		if key != "" {
			data = url.Values{key: []string{value}}.Encode()
		}

		req := httptest.NewRequest("POST", "/ping", strings.NewReader(data))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		return req
	}

	newMultipartRequest := func(key, value string) (*http.Request, error) {
		contentType := "multipart/form-data"

		var c int64
		var data bytes.Buffer
		if key != "" {
			w := multipart.NewWriter(&data)

			fw, err := w.CreateFormField(key)
			if err != nil {
				return nil, fmt.Errorf("error creating field: %s: %w", key, err)
			}

			c, err = io.Copy(fw, strings.NewReader(value))
			if err != nil {
				return nil, fmt.Errorf("error copying value for field: %s: %w", key, err)
			}

			err = w.Close()
			if err != nil {
				return nil, fmt.Errorf("error closing writer: for field: %s: %w", key, err)
			}

			contentType = w.FormDataContentType()
		}

		req := httptest.NewRequest("POST", "/ping", &data)
		req.Header.Set("Content-Type", contentType)

		fmt.Printf("request: %v: %v", c, req)

		return req, nil
	}

	tests := []struct {
		name     string
		key      string
		value    string
		form     string
		required bool
		want     string
		wantErr  bool
	}{
		{" required present", "foo", "123", "foo", true, "123", false},
		{" required missing", "", "", "foo", true, "", true},
		{" not required present", "foo", "123", "foo", false, "123", false},
		{" not required missing", "", "", "foo", false, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify that GetString returns data from each of the FormValue sources:
			//  1. application/x-www-form-urlencoded form body (POST, PUT, PATCH only)
			//  2. query parameters (always)
			//  3. multipart/form-data form body (always)
			var params string
			if tt.key != "" {
				params = fmt.Sprintf("?%s=%s", tt.key, tt.value)
			}

			r := httptest.NewRequest("POST", "/BAR"+params, nil)

			got, err := GetString(r.FormValue, tt.form, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")

			r = newFormRequest(tt.key, tt.value)

			got, err = GetString(r.FormValue, tt.form, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")

			r, err = newMultipartRequest(tt.key, tt.value)
			assert.NoError(t, err)

			got, err = GetString(r.FormValue, tt.form, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
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
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=123", nil), "foo", true, 123, false},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, 0, true},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, 0, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInt32(tt.r.FormValue, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
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
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=123", nil), "foo", true, 123, false},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, 0, true},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, 0, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, 0, true},
		{"max", httptest.NewRequest("GET", "/BAR?foo=9223372036854775807", nil), "foo", true, 9223372036854775807, false},
		{"over max", httptest.NewRequest("GET", "/BAR?foo=19223372036854775807", nil), "foo", true, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInt(tt.r.FormValue, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
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
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=123", nil), "foo", true, 123, false},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, 0, true},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, 0, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, 0, true},
		{"max", httptest.NewRequest("GET", "/BAR?foo=9223372036854775807", nil), "foo", true, 9223372036854775807, false},
		{"over max", httptest.NewRequest("GET", "/BAR?foo=19223372036854775807", nil), "foo", true, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInt64(tt.r.FormValue, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
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
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=true", nil), "foo", true, true, false},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, false, true},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, false, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBool(tt.r.FormValue, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
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
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=123", nil), "foo", true, []int32{123}, false},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, nil, true},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, nil, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, nil, true},
		{"large", httptest.NewRequest("GET", "/BAR?foo=1,2,3,4", nil), "foo", true, []int32{1, 2, 3, 4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInt32Array(tt.r.FormValue, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
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
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=2006-01-02T15:04:05Z", nil), "foo", true, time.Date(2006, 0o1, 0o2, 15, 4, 5, 0, time.UTC), false},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, time.Time{}, true},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, time.Time{}, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=200601021504050700", nil), "foo", true, time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTime(tt.r.FormValue, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
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
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=48ab873f-d4fc-4e2b-bf92-9440e431ff54", nil), "foo", true, testUUID, false},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, uuid.UUID{}, true},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, uuid.UUID{}, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, uuid.UUID{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUUID(tt.r.FormValue, tt.param, tt.required)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
		})
	}
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
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=bbb", nil), "foo", true, testEnumB, false},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, testEnumUnknown, true},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, testEnumUnknown, false},
		{"bad value", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, testEnumUnknown, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEnum(tt.r.FormValue, tt.param, tt.required, NewTestEnum)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
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
	}{
		{"simple", httptest.NewRequest("GET", "/BAR?foo=bbb", nil), "foo", true, []TestEnum{testEnumB}, false},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, nil, true},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, nil, false},
		{"bad value", httptest.NewRequest("GET", "/BAR?foo=a123", nil), "foo", true, []TestEnum{testEnumUnknown}, false},
		{"all", httptest.NewRequest("GET", "/BAR?foo=aaa,bbb,ccc", nil), "foo", true, []TestEnum{testEnumA, testEnumB, testEnumC}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEnumArray(tt.r.FormValue, tt.param, tt.required, NewTestEnum)

			assert.Equal(t, tt.wantErr, err != nil, "error")
			assert.Equal(t, tt.want, got, "value")
		})
	}
}

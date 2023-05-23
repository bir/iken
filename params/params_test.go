package params

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

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
			got, err := GetInt32(tt.r, tt.param, tt.required)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				if tt.want != 0 {
					t.Errorf("GetInt32() got = %v, want %v", got, tt.want)
				}
			} else if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("GetInt32() got = %v, want %v", *got, tt.want)
			}
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
			got, err := GetInt(tt.r, tt.param, tt.required)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				if tt.want != 0 {
					t.Errorf("GetInt() got = %v, want %v", got, tt.want)
				}
			} else if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("GetInt() got = %v, want %v", *got, tt.want)
			}
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
			got, err := GetInt32Array(tt.r, tt.param, tt.required)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInt32Array() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInt32Array() got = %v, want %v", got, tt.want)
			}
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
		{"simple", httptest.NewRequest("GET", "/BAR?foo=2006-01-02T15:04:05Z", nil), "foo", true, time.Date(2006, 01, 02, 15, 4, 5, 0, time.UTC), false},
		{"required missing", httptest.NewRequest("GET", "/BAR", nil), "foo", true, time.Time{}, true},
		{"not required missing", httptest.NewRequest("GET", "/BAR?", nil), "foo", false, time.Time{}, false},
		{"bad format", httptest.NewRequest("GET", "/BAR?foo=200601021504050700", nil), "foo", true, time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTime(tt.r, tt.param, tt.required)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

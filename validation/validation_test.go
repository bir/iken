package validation_test

import (
	"testing"

	"github.com/bir/iken/validation"
)

func TestErrors_Add(t *testing.T) {
	tests := []struct {
		name  string
		ee    validation.Errors
		field string
		msg   string
		want  string
	}{
		{"basic", *(&validation.Errors{}).Add("a", "b"), "test", "bad", "a: b; test: bad."},
		{"existing", *(&validation.Errors{}).Add("a", "b"), "a", "x", "a: b, x."},
		{"basic nil", nil, "test", "bad", "test: bad."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ee.Add(tt.field, tt.msg); got.Error() != tt.want {
				t.Errorf("Add() = `%v`, want `%v`", got.Error(), tt.want)
			}
		})
	}
}

func TestErrors_Error(t *testing.T) {
	var ee validation.Errors
	if ee.Error() != "" {
		t.Errorf("Error() = `%v`, want ``", ee.Error())
	}
}

func TestErrors_GetErr(t *testing.T) {
	var ee validation.Errors
	if ee.GetErr() != nil {
		t.Errorf("GetErr() = `%v`, want `nil`", ee.GetErr())
	}
	_ = ee.Add("a","b")

	if ee.GetErr() == nil {
		t.Errorf("GetErr() = `nil`, want `a: b.`")
	}
}

func TestErrors_New(t *testing.T) {
	err := validation.New("a","b")
	if err == nil {
		t.Errorf("New() = `nil`, want `a: b.`")
	}

	if err.Error() != "a: b." {
		t.Errorf("New() = `%v`, want `a: b.`", err)
	}

}

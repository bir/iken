package httputil_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bir/iken/arrays"
	"github.com/bir/iken/httputil"
)

type (
	AuthorizeFunc  = httputil.AuthorizeFunc[string]
	SecurityGroup  = httputil.SecurityGroup[string]
	SecurityGroups = httputil.SecurityGroups[string]
)

func authenticate(r *http.Request) (string, error) {
	hdr := r.Header.Get("Authorization")
	switch hdr {
	case "tokenForA":
		return "A", nil
	case "tokenForB":
		return "B", nil
	}

	return "", errors.New("missing")
}

func authorize(ctx context.Context, user string, scopes []string) error {
	if arrays.Contains(user, scopes) {
		return nil
	}

	return errors.New("bad")
}

func TestAuthCheck_Auth(t *testing.T) {

	type testCase struct {
		name      string
		authorize AuthorizeFunc
		scopes    []string
		hdr       string
		want      string
		wantErr   bool
	}

	tests := []testCase{
		{"unknown", authorize, nil, "blah", "", true},
		{"A", authorize, nil, "tokenForA", "A", false},
		{"B", authorize, nil, "tokenForB", "B", false},
		{"scoped A", authorize, []string{"A"}, "tokenForA", "A", false},
		{"scoped A missing", authorize, []string{"B"}, "tokenForA", "", true},
		{"scoped missing authorize", nil, []string{"A"}, "tokenForA", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chk := httputil.NewAuthCheck(
				authenticate,
				tt.authorize,
				tt.scopes...)
			r := httptest.NewRequest("FOO", "/asdf", nil)
			r.Header.Set("Authorization", tt.hdr)
			got, err := chk.Auth(r)

			if !tt.wantErr {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecurityGroup_Auth(t *testing.T) {
	type testCase struct {
		name    string
		s       SecurityGroup
		hdr     string
		want    string
		wantErr bool
	}

	chkA := httputil.NewAuthCheck(authenticate, authorize, "A")
	chkNoScope := httputil.NewAuthCheck(authenticate, nil)
	chkB := httputil.NewAuthCheck(authenticate, authorize, "B")

	tests := []testCase{
		{"basic A", SecurityGroup{chkA, chkNoScope}, "tokenForA", "A", false},
		{"impossible A", SecurityGroup{chkA, chkB}, "tokenForA", "", true},
		{"missing scope", SecurityGroup{chkA}, "tokenForB", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)
			r.Header.Set("Authorization", tt.hdr)

			got, err := tt.s.Auth(r)
			if !tt.wantErr {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSecurityGroups_Auth(t *testing.T) {
	type testCase[T any] struct {
		name    string
		s       SecurityGroups
		hdr     string
		want    T
		wantErr bool
	}

	chkA := httputil.NewAuthCheck(authenticate, authorize, "A")
	chkNoScope := httputil.NewAuthCheck(authenticate, nil)
	chkB := httputil.NewAuthCheck(authenticate, authorize, "B")

	tests := []testCase[string]{
		{"A||B", SecurityGroups{SecurityGroup{chkA}, SecurityGroup{chkB}}, "tokenForA", "A", false},
		{"A&&true||B", SecurityGroups{SecurityGroup{chkA, chkNoScope}, SecurityGroup{chkB}}, "tokenForB", "B", false},
		{"none", SecurityGroups{}, "tokenForB", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)
			r.Header.Set("Authorization", tt.hdr)

			got, err := tt.s.Auth(r)
			if !tt.wantErr {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

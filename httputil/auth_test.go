package httputil_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func authorize(_ context.Context, user string, scopes []string) error {
	if slices.Contains(scopes, user) {
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

func strAuth(_ context.Context, token string) (string, error) {
	if token == "" {
		return "", errors.New("unreachable")
	}

	if token == "good" {
		return "good", nil
	}

	if token == "bad" {
		return "", errors.New("bad")
	}

	return "", errors.New("badder")
}

func TestHeaderAuth(t *testing.T) {
	type testCase[T any] struct {
		name string
		key  string
		val  string
		fn   httputil.TokenAuthenticatorFunc[T]
		want string
		err  error
	}
	tests := []testCase[string]{
		{"Empty", "Missing", "", strAuth, "", httputil.ErrUnauthorized},
		{"Good", "Authorization", "good", strAuth, "good", nil},
		{"Bad", "Authorization", "bad", strAuth, "", errors.New("bad")},
		{"other", "Authorization", "other", strAuth, "", errors.New("badder")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)
			r.Header.Set(tt.key, tt.val)

			got, err := httputil.HeaderAuth(tt.key, tt.fn)(r)
			if tt.err != nil {
				assert.Equal(t, tt.err, err)
			}

			assert.Equalf(t, tt.want, got, "HeaderAuth(%v, %v)", tt.key, tt.fn)
		})
	}
}

func TestQueryAuth(t *testing.T) {
	type testCase[T any] struct {
		name string
		key  string
		val  string
		fn   httputil.TokenAuthenticatorFunc[T]
		want string
		err  error
	}
	tests := []testCase[string]{
		{"Empty", "Missing", "", strAuth, "", httputil.ErrUnauthorized},
		{"Good", "Authorization", "good", strAuth, "good", nil},
		{"Bad", "Authorization", "bad", strAuth, "", errors.New("bad")},
		{"other", "Authorization", "other", strAuth, "", errors.New("badder")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", fmt.Sprintf("/asdf?%s=%s", tt.key, tt.val), nil)

			got, err := httputil.QueryAuth(tt.key, tt.fn)(r)
			if tt.err != nil {
				assert.Equal(t, tt.err, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCookieAuth(t *testing.T) {
	type testCase[T any] struct {
		name string
		key  string
		val  string
		fn   httputil.TokenAuthenticatorFunc[T]
		want string
		err  error
	}
	tests := []testCase[string]{
		{"Empty", "Missing", "", strAuth, "", httputil.ErrUnauthorized},
		{"Good", "Authorization", "good", strAuth, "good", nil},
		{"Bad", "Authorization", "bad", strAuth, "", errors.New("bad")},
		{"other", "Authorization", "other", strAuth, "", errors.New("badder")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)
			r.AddCookie(&http.Cookie{
				Name:  tt.key,
				Value: tt.val,
			})

			got, err := httputil.CookieAuth(tt.key, tt.fn)(r)
			if tt.err != nil {
				assert.Equal(t, tt.err, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBearerAuth(t *testing.T) {
	type testCase[T any] struct {
		name string
		key  string
		val  string
		fn   httputil.TokenAuthenticatorFunc[T]
		want string
		err  error
	}
	tests := []testCase[string]{
		{"Empty", "Missing", "", strAuth, "", httputil.ErrUnauthorized},
		{"Empty Bearer", "Missing", "Bearer ", strAuth, "", httputil.ErrUnauthorized},
		{"Good optional bearer", "Authorization", "good", strAuth, "good", nil},
		{"Good", "Authorization", "Bearer good", strAuth, "good", nil},
		{"Bad", "Authorization", "Bearer bad", strAuth, "", errors.New("bad")},
		{"other", "Authorization", "Bearer other", strAuth, "", errors.New("badder")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)
			r.Header.Set(tt.key, tt.val)

			got, err := httputil.BearerAuth(tt.key, tt.fn)(r)
			if tt.err != nil {
				assert.Equal(t, err, tt.err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

var ErrBad = errors.New("bad")

func basicAuth(_ context.Context, user, _ string) (string, error) {
	if user == "" {
		return "", errors.New("unreachable")
	}

	if user == "good" {
		return "good", nil
	}

	if user == "bad" {
		return "", ErrBad
	}

	return "", errors.New("badder")
}

func TestBasicAuth(t *testing.T) {
	type testCase[T any] struct {
		name string
		key  string
		val  string
		fn   httputil.BasicAuthenticatorFunc[T]
		want string
		err  error
	}
	tests := []testCase[string]{
		{"Empty", "Missing", "", basicAuth, "", httputil.ErrBasicAuthenticate},
		{"Invalid base64", "Authorization", "Basic other", basicAuth, "", httputil.ErrBasicAuthenticate},
		{"No split", "Authorization", "Basic Z29vZGJhZA==", basicAuth, "", httputil.ErrUnauthorized},
		{"Good", "Authorization", "Basic Z29vZDpwYXNz", basicAuth, "good", nil},
		{"Good Proxy", "Proxy-Authorization", "Basic Z29vZDpwYXNz", basicAuth, "good", nil},
		{"Bad", "Authorization", "Basic YmFkOnBhc3M=", basicAuth, "", ErrBad},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)
			r.Header.Set(tt.key, tt.val)

			got, err := httputil.BasicAuth(tt.fn)(r)
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func strClientAuth(_ context.Context, user string) (string, error) {
	switch user {
	case "good":
		return "token", nil
	case "bad":
		return "", ErrBad
	}

	return "", httputil.ErrCannotAuthenticate
}

func failClientAuth(_ context.Context, _ string) (string, error) {
	return "", ErrBad
}

func onlyUser(name string) httputil.ClientTokenAuthenticatorFunc[string] {
	return func(_ context.Context, user string) (string, error) {
		if user == name {
			return "token", nil
		}

		return "", httputil.ErrCannotAuthenticate
	}
}

func basicClientAuth(_ context.Context, user string) (string, string, error) {
	switch user {
	case "good":
		return "u", "p", nil
	case "bad":
		return "", "", ErrBad
	}

	return "", "", httputil.ErrCannotAuthenticate
}

func cookieClientAuth(_ context.Context, user string) (*http.Cookie, error) {
	switch user {
	case "good":
		return &http.Cookie{Name: "session", Value: "token"}, nil
	case "nil":
		return nil, nil
	case "bad":
		return nil, ErrBad
	}

	return nil, httputil.ErrCannotAuthenticate
}

func TestHeaderClientAuth(t *testing.T) {
	type testCase struct {
		name string
		key  string
		user string
		want string
		err  error
	}
	tests := []testCase{
		{"Good", "X-Auth", "good", "token", nil},
		{"Bad", "X-Auth", "bad", "", ErrBad},
		{"Couldnt", "X-Auth", "other", "", httputil.ErrCannotAuthenticate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)
			inner := &http.Client{}

			got, err := httputil.HeaderClientAuth(tt.key, strClientAuth)(r, inner, tt.user)
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				assert.Same(t, inner, got)
			}

			assert.Equal(t, tt.want, r.Header.Get(tt.key))
		})
	}
}

func TestBearerClientAuth(t *testing.T) {
	type testCase struct {
		name string
		key  string
		user string
		want string
		err  error
	}
	tests := []testCase{
		{"Good", "Authorization", "good", "Bearer token", nil},
		{"Bad", "Authorization", "bad", "", ErrBad},
		{"Couldnt", "Authorization", "other", "", httputil.ErrCannotAuthenticate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)
			inner := &http.Client{}

			got, err := httputil.BearerClientAuth(tt.key, strClientAuth)(r, inner, tt.user)
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				assert.Same(t, inner, got)
			}

			assert.Equal(t, tt.want, r.Header.Get(tt.key))
		})
	}
}

func TestQueryClientAuth(t *testing.T) {
	type testCase struct {
		name string
		key  string
		user string
		want string
		err  error
	}
	tests := []testCase{
		{"Good", "auth", "good", "token", nil},
		{"Bad", "auth", "bad", "", ErrBad},
		{"Couldnt", "auth", "other", "", httputil.ErrCannotAuthenticate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)
			inner := &http.Client{}

			got, err := httputil.QueryClientAuth(tt.key, strClientAuth)(r, inner, tt.user)
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				assert.Same(t, inner, got)
			}

			assert.Equal(t, tt.want, r.URL.Query().Get(tt.key))
		})
	}
}

func TestBasicClientAuth(t *testing.T) {
	type testCase struct {
		name     string
		user     string
		wantUser string
		wantPass string
		err      error
	}
	tests := []testCase{
		{"Good", "good", "u", "p", nil},
		{"Bad", "bad", "", "", ErrBad},
		{"Couldnt", "other", "", "", httputil.ErrCannotAuthenticate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)
			inner := &http.Client{}

			got, err := httputil.BasicClientAuth(basicClientAuth)(r, inner, tt.user)
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				assert.Same(t, inner, got)
			}

			u, p, ok := r.BasicAuth()
			if tt.err == nil {
				assert.True(t, ok)
				assert.Equal(t, tt.wantUser, u)
				assert.Equal(t, tt.wantPass, p)
			} else {
				assert.False(t, ok)
			}
		})
	}
}

func TestCookieClientAuth(t *testing.T) {
	type testCase struct {
		name    string
		user    string
		want    string
		wantErr error
	}
	tests := []testCase{
		{"Good", "good", "token", nil},
		{"Nil", "nil", "", nil},
		{"Bad", "bad", "", ErrBad},
		{"Couldnt", "other", "", httputil.ErrCannotAuthenticate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)
			inner := &http.Client{}

			got, err := httputil.CookieClientAuth(cookieClientAuth)(r, inner, tt.user)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Same(t, inner, got)
			}

			cookie, cookieErr := r.Cookie("session")
			if tt.want != "" {
				require.NoError(t, cookieErr)
				assert.Equal(t, tt.want, cookie.Value)
			} else {
				assert.ErrorIs(t, cookieErr, http.ErrNoCookie)
			}
		})
	}
}

func TestWrapClientAuth(t *testing.T) {
	inner := &http.Client{}
	wrapped := &http.Client{}

	fn := func(_ context.Context, in *http.Client, user string) (*http.Client, error) {
		switch user {
		case "good":
			return wrapped, nil
		case "bad":
			return nil, ErrBad
		}

		return in, httputil.ErrCannotAuthenticate
	}

	type testCase struct {
		name string
		user string
		want *http.Client
		err  error
	}
	tests := []testCase{
		{"good", "good", wrapped, nil},
		{"bad", "bad", inner, ErrBad},
		{"couldnt", "other", inner, httputil.ErrCannotAuthenticate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)

			got, err := httputil.WrapClientAuth(fn)(r, inner, tt.user)
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				assert.Same(t, tt.want, got)
			}
		})
	}
}

func TestClientSecurityGroup_Auth(t *testing.T) {
	inner := &http.Client{}

	type testCase struct {
		name  string
		g     httputil.ClientSecurityGroup[string]
		user  string
		wantA string
		wantB string
		err   error
	}
	tests := []testCase{
		{
			"all succeed",
			httputil.ClientSecurityGroup[string]{
				httputil.HeaderClientAuth("X-A", strClientAuth),
				httputil.HeaderClientAuth("X-B", strClientAuth),
			},
			"good", "token", "token", nil,
		},
		{
			"failure returns inner and leaves request unmodified",
			httputil.ClientSecurityGroup[string]{
				httputil.HeaderClientAuth("X-A", strClientAuth),
				httputil.HeaderClientAuth("X-B", failClientAuth),
			},
			"good", "", "", ErrBad,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)

			got, err := tt.g.Auth(r, inner, tt.user)
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				assert.Same(t, inner, got)
			}

			assert.Equal(t, tt.wantA, r.Header.Get("X-A"))
			assert.Equal(t, tt.wantB, r.Header.Get("X-B"))
		})
	}
}

func TestClientSecurityGroups_Auth(t *testing.T) {
	inner := &http.Client{}

	groups := httputil.ClientSecurityGroups[string]{
		{httputil.HeaderClientAuth("X-A", onlyUser("a"))},
		{httputil.HeaderClientAuth("X-B", onlyUser("b"))},
	}
	failing := httputil.ClientSecurityGroups[string]{
		{httputil.HeaderClientAuth("X-A", failClientAuth)},
		{httputil.HeaderClientAuth("X-B", onlyUser("b"))},
	}

	type testCase struct {
		name  string
		s     httputil.ClientSecurityGroups[string]
		user  string
		wantA string
		wantB string
		err   error
	}
	tests := []testCase{
		{"first group succeeds", groups, "a", "token", "", nil},
		{"falls through to a later group", groups, "b", "", "token", nil},
		{"none can authenticate", groups, "c", "", "", httputil.ErrCannotAuthenticate},
		{"true failure stops iteration", failing, "b", "", "", ErrBad},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("FOO", "/asdf", nil)

			got, err := tt.s.Auth(r, inner, tt.user)
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				assert.Same(t, inner, got)
			}

			assert.Equal(t, tt.wantA, r.Header.Get("X-A"))
			assert.Equal(t, tt.wantB, r.Header.Get("X-B"))
		})
	}
}

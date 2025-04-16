package httputil

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// AuthError encompasses Authentication and Authorization errors.
type AuthError string

func (e AuthError) Error() string {
	return string(e)
}

const (
	// ErrUnauthorized represents failure when authenticating a request.
	ErrUnauthorized = AuthError("Unauthorized")

	// ErrForbidden represents failure when authorizing a request.
	ErrForbidden = AuthError("Forbidden")

	// ErrBasicAuthenticate is designed to trigger an HTTP Basic Auth challenge.
	ErrBasicAuthenticate = AuthError("ErrWWWAuthenticate")

	// ErrMissingAuthorizer is caused by internal configuration errors when evaluating authorization.
	ErrMissingAuthorizer = AuthError("missing authenticator")
)

// AuthenticateFunc is the signature of a function used to authenticate an HTTP request.
// Given a request, it returns the authenticated user.  If unable to authenticate the
// request it returns an error.
type AuthenticateFunc[T any] func(r *http.Request) (T, error)

// TokenAuthenticatorFunc is the signature of a function used to authenticate a request given just the token.
// Given a request, it returns the authenticated user.  If unable to authenticate the
// request it returns an error.
type TokenAuthenticatorFunc[T any] func(ctx context.Context, token string) (T, error)

// AuthorizeFunc is the signature of a function used to authorize a request.  If unable
// to authorize the user it returns an error.
type AuthorizeFunc[T any] func(ctx context.Context, user T, scopes []string) error

// AuthCheck combines the authenticator with an authorizer and set of scopes.  This is generally attached
// to each end point for auth.
type AuthCheck[T any] struct {
	authenticate AuthenticateFunc[T]
	authorize    AuthorizeFunc[T]
	scopes       []string
}

func HeaderAuth[T any](key string, fn TokenAuthenticatorFunc[T]) AuthenticateFunc[T] {
	return func(r *http.Request) (T, error) {
		var empty T

		token := r.Header.Get(key)
		if token == "" {
			return empty, ErrUnauthorized
		}

		return fn(r.Context(), token)
	}
}

const bearerAuthPrefix = "Bearer "

func BearerAuth[T any](key string, tokenAuth TokenAuthenticatorFunc[T]) AuthenticateFunc[T] {
	return func(r *http.Request) (T, error) {
		var empty T

		token := r.Header.Get(key)

		if strings.HasPrefix(token, bearerAuthPrefix) {
			token = token[7:]
		}

		if token == "" {
			return empty, ErrUnauthorized
		}

		return tokenAuth(r.Context(), token)
	}
}

func QueryAuth[T any](key string, fn TokenAuthenticatorFunc[T]) AuthenticateFunc[T] {
	return func(r *http.Request) (T, error) {
		var empty T

		token := r.URL.Query().Get(key)
		if token == "" {
			return empty, ErrUnauthorized
		}

		return fn(r.Context(), token)
	}
}

func CookieAuth[T any](key string, fn TokenAuthenticatorFunc[T]) AuthenticateFunc[T] {
	return func(r *http.Request) (T, error) {
		var empty T

		cookie, err := r.Cookie(key)
		if err != nil || cookie == nil || len(cookie.Value) == 0 {
			return empty, ErrUnauthorized
		}

		return fn(r.Context(), cookie.Value)
	}
}

func NewAuthCheck[T any](authenticate AuthenticateFunc[T], authorize AuthorizeFunc[T], scopes ...string) AuthCheck[T] {
	return AuthCheck[T]{
		authenticate: authenticate,
		authorize:    authorize,
		scopes:       scopes,
	}
}

func (a AuthCheck[T]) Auth(r *http.Request) (T, error) {
	var empty T

	user, err := a.authenticate(r)
	if err != nil {
		return empty, fmt.Errorf("%w:%w", ErrUnauthorized, err)
	}

	if len(a.scopes) > 0 {
		if a.authorize == nil {
			return empty, ErrMissingAuthorizer
		}

		err = a.authorize(r.Context(), user, a.scopes)
		if err != nil {
			return empty, fmt.Errorf("%w:%w", ErrForbidden, err)
		}
	}

	return user, nil
}

// SecurityGroup are valid if all are valid.
type SecurityGroup[T any] []AuthCheck[T]

// Auth returns a user if all AuthChecks are successful.
func (s SecurityGroup[T]) Auth(r *http.Request) (T, error) {
	var (
		user T
		err  error
	)

	for _, check := range s {
		user, err = check.Auth(r)
		if err != nil {
			var empty T

			return empty, err
		}
	}

	return user, nil
}

// SecurityGroups are valid if ANY group is valid.
type SecurityGroups[T any] []SecurityGroup[T]

// Auth returns a user if any of the group checks is successful.
func (s SecurityGroups[T]) Auth(r *http.Request) (T, error) {
	var (
		err  error
		user T
	)

	for _, group := range s {
		user, err = group.Auth(r)
		if err == nil {
			return user, nil
		}
	}

	return user, err
}

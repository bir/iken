package httputil

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
)

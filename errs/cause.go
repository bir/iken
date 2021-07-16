package errs

// Cause first checks for `Unwrap` returning the result if available,
// otherwise returns the result of `Cause`.
// If neither are available returns nil.
//
// This supports both Go 1.13+ style `UnWrap` and pkg.errors style
// `Cause` chaining.
func Cause(err error) error {
	u, ok := err.(interface{ Unwrap() error }) //nolint:errorlint // false positive
	if ok {
		return u.Unwrap() //nolint:wrapcheck // defeats the whole point
	}

	c, ok := err.(interface{ Cause() error }) //nolint:errorlint // false positive
	if ok {
		return c.Cause() //nolint:wrapcheck // defeats the whole point
	}

	return nil
}

// RootCause follows the error chain until Cause() and UnWrap() are not
// available or return nil.
func RootCause(err error) error {
	for err != nil {
		errCause := Cause(err)
		if errCause == nil {
			return err
		}

		err = errCause
	}

	return err
}

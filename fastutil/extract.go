package fastutil

import (
	"strconv"

	"github.com/valyala/fasthttp"
)

// ErrRequired is returned by request param parser helpers when the requested params are not provided.
const ErrRequired = ExtractError("required")

// ExtractError is the static type for param extraction errors.
type ExtractError string

func (e ExtractError) Error() string {
	return string(e)
}

// QueryString returns the string from the query identified by key, ErrRequired otherwise.
func QueryString(ctx *fasthttp.RequestCtx, key string) (string, error) {
	b := ctx.QueryArgs().Peek(key)
	if b == nil {
		return "", ErrRequired
	}

	s := string(b)

	return s, nil
}

// QueryStringOptional returns the *string from the query identified by key, nil otherwise.
func QueryStringOptional(ctx *fasthttp.RequestCtx, key string) *string {
	b := ctx.QueryArgs().Peek(key)
	if b == nil {
		return nil
	}

	s := string(b)

	return &s
}

// QueryStringsOptional returns the array of string from the query identified by key.  For example
// ?test=1&test=2 will return 1,2.
func QueryStringsOptional(ctx *fasthttp.RequestCtx, key string) []string {
	pp := ctx.QueryArgs().PeekMulti(key)
	out := make([]string, len(pp))

	for i := range pp {
		out[i] = string(pp[i])
	}

	return out
}

// QueryInt32Optional returns the *int32 from the query identified by key, nil otherwise.
func QueryInt32Optional(ctx *fasthttp.RequestCtx, key string) *int32 {
	out, err := QueryInt32(ctx, key)
	if err != nil {
		return nil
	}

	return &out
}

// QueryInt32 returns the *int32 from the query identified by key, err otherwise.
func QueryInt32(ctx *fasthttp.RequestCtx, key string) (int32, error) {
	out, err := ctx.QueryArgs().GetUint(key)
	if err != nil {
		return 0, err
	}

	return int32(out), nil
}

// QueryInt64Optional returns the *int64 from the query identified by key, nil otherwise.
func QueryInt64Optional(ctx *fasthttp.RequestCtx, key string) *int64 {
	out, err := QueryInt64(ctx, key)
	if err != nil {
		return nil
	}

	return &out
}

// QueryInt64 returns the *int64 from the query identified by key, err otherwise.
func QueryInt64(ctx *fasthttp.RequestCtx, key string) (int64, error) {
	out, err := ctx.QueryArgs().GetUint(key)
	if err != nil {
		return 0, err
	}

	return int64(out), nil
}

// QueryBool returns the bool from the query identified by key, err otherwise.
//
// true is returned for "1", "t", "T", "true", "TRUE", "True", "y", "yes", "Y", "YES", "Yes"
// Error is always nil, it return an error so that the contract for usage is consistent for
// all helper functions.
func QueryBool(ctx *fasthttp.RequestCtx, key string) (bool, error) {
	return ctx.QueryArgs().GetBool(key), nil
}

// PathInt64 returns the int64 from the Path Param identified by key, err otherwise.
func PathInt64(ctx *fasthttp.RequestCtx, key string) (int64, error) {
	s, ok := ctx.UserValue(key).(string)
	if !ok || len(s) == 0 {
		return 0, ErrRequired
	}

	return strconv.ParseInt(s, 0, 64)
}

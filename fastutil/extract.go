package fastutil

import (
	"strconv"

	"github.com/valyala/fasthttp"
)

const ErrRequired = ExtractError("required")

type ExtractError string

func (e ExtractError) Error() string {
	return string(e)
}

func QueryString(ctx *fasthttp.RequestCtx, key string) (string, error) {
	b := ctx.QueryArgs().Peek(key)
	if b == nil {
		return "", ErrRequired
	}

	s := string(b)

	return s, nil
}

func QueryStringOptional(ctx *fasthttp.RequestCtx, key string) *string {
	b := ctx.QueryArgs().Peek(key)
	if b == nil {
		return nil
	}

	s := string(b)

	return &s
}

func QueryStringsOptional(ctx *fasthttp.RequestCtx, key string) []string {
	pp := ctx.QueryArgs().PeekMulti(key)
	out := make([]string, len(pp))

	for i := range pp {
		out[i] = string(pp[i])
	}

	return out
}

func QueryInt32Optional(ctx *fasthttp.RequestCtx, key string) *int32 {
	out, err := QueryInt32(ctx, key)
	if err != nil {
		return nil
	}

	return &out
}

func QueryInt32(ctx *fasthttp.RequestCtx, key string) (int32, error) {
	out, err := ctx.QueryArgs().GetUint(key)
	if err != nil {
		return 0, err
	}

	return int32(out), nil
}

func QueryInt64Optional(ctx *fasthttp.RequestCtx, key string) *int64 {
	out, err := QueryInt64(ctx, key)
	if err != nil {
		return nil
	}

	return &out
}

func QueryInt64(ctx *fasthttp.RequestCtx, key string) (int64, error) {
	out, err := ctx.QueryArgs().GetUint(key)
	if err != nil {
		return 0, err
	}

	return int64(out), nil
}

func QueryBool(ctx *fasthttp.RequestCtx, key string) (bool, error) {
	return ctx.QueryArgs().GetBool(key), nil
}

func PathInt64(ctx *fasthttp.RequestCtx, key string) (int64, error) {
	s, ok := ctx.UserValue(key).(string)
	if !ok || len(s) == 0 {
		return 0, ErrRequired
	}

	return strconv.ParseInt(s, 0, 64)
}

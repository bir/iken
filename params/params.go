package params

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrUnknownParamSource = errors.New("param source unknown")
)

type ParamSource int

const (
	ParamPath ParamSource = iota
	ParamQuery
	ParamHeader
	ParamCookie
)

func GetString(r *http.Request, name string, required bool) (string, bool, error) {
	param := r.PathValue(name)

	if param == "" {
		param = r.URL.Query().Get(name)
	}

	// fallback to a header lookup
	if param == "" {
		param = r.Header.Get(name)
	}

	if required && len(param) == 0 {
		return "", false, fmt.Errorf("%s: %w", name, ErrNotFound)
	}

	return param, param != "", nil
}

func GetStringFrom(r *http.Request, name string, source ParamSource, required bool) (string, bool, error) {
	var param string

	switch source {
	case ParamPath:
		param = r.PathValue(name)
	case ParamQuery:
		param = r.URL.Query().Get(name)
	case ParamHeader:
		param = r.Header.Get(name)
	case ParamCookie:
		cookie, err := r.Cookie(name)
		// only error is cookie not found, so leave param blank in that case.
		if err == nil && cookie != nil {
			param = cookie.Value
		}
	default:
		return "", false, fmt.Errorf("%d: %w", source, ErrUnknownParamSource)
	}

	if required && len(param) == 0 {
		return "", false, fmt.Errorf("%s: %w", name, ErrNotFound)
	}

	return param, param != "", nil
}

func GetInt32(r *http.Request, name string, required bool) (int32, bool, error) {
	s, ok, err := GetString(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt32(s)
}

func GetInt32From(r *http.Request, name string, source ParamSource, required bool) (int32, bool, error) {
	s, ok, err := GetStringFrom(r, name, source, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt32(s)
}

func convertInt32(s string) (int32, bool, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, false, fmt.Errorf("invalid int32: %w", err)
	}

	return int32(i), true, nil
}

func GetInt64(r *http.Request, name string, required bool) (int64, bool, error) {
	s, ok, err := GetString(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt64(s)
}

func GetInt64From(r *http.Request, name string, source ParamSource, required bool) (int64, bool, error) {
	s, ok, err := GetStringFrom(r, name, source, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt64(s)
}

func convertInt64(s string) (int64, bool, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("invalid int32: %w", err)
	}

	return i, true, nil
}

func GetBool(r *http.Request, name string, required bool) (bool, bool, error) {
	s, ok, err := GetString(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return false, false, err
	}

	return convertBool(s)
}

func GetBoolFrom(r *http.Request, name string, source ParamSource, required bool) (bool, bool, error) {
	s, ok, err := GetStringFrom(r, name, source, required)
	if err != nil || len(s) == 0 || !ok {
		return false, false, err
	}

	return convertBool(s)
}

func convertBool(s string) (bool, bool, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, false, fmt.Errorf("invalid bool: %w", err)
	}

	return b, true, nil
}

func GetInt(r *http.Request, name string, required bool) (int, bool, error) {
	s, ok, err := GetString(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt(s)
}

func GetIntFrom(r *http.Request, name string, source ParamSource, required bool) (int, bool, error) {
	s, ok, err := GetStringFrom(r, name, source, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt(s)
}

func convertInt(s string) (int, bool, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false, fmt.Errorf("invalid int: %w", err)
	}

	return i, true, nil
}

func GetTime(r *http.Request, name string, required bool) (time.Time, bool, error) {
	s, ok, err := GetString(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return time.Time{}, false, err
	}

	return convertTime(s)
}

func GetTimeFrom(r *http.Request, name string, source ParamSource, required bool) (time.Time, bool, error) {
	s, ok, err := GetStringFrom(r, name, source, required)
	if err != nil || len(s) == 0 || !ok {
		return time.Time{}, false, err
	}

	return convertTime(s)
}

func convertTime(s string) (time.Time, bool, error) {
	timestamp, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("invalid RFC3339 date: %w", err)
	}

	return timestamp, true, nil
}

func GetUUID(r *http.Request, name string, required bool) (uuid.UUID, bool, error) {
	s, ok, err := GetString(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return uuid.UUID{}, false, err
	}

	return convertUUID(s)
}

func GetUUIDFrom(r *http.Request, name string, source ParamSource, required bool) (uuid.UUID, bool, error) {
	s, ok, err := GetStringFrom(r, name, source, required)
	if err != nil || len(s) == 0 || !ok {
		return uuid.UUID{}, false, err
	}

	return convertUUID(s)
}

func convertUUID(s string) (uuid.UUID, bool, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, false, fmt.Errorf("invalid uuid: %w", err)
	}

	return id, true, nil
}

func GetStringArray(r *http.Request, name string, required bool) ([]string, bool, error) {
	s, ok, err := GetString(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return nil, false, err
	}

	return strings.Split(s, ","), true, nil
}

func GetStringArrayFrom(r *http.Request, name string, source ParamSource, required bool) ([]string, bool, error) {
	s, ok, err := GetStringFrom(r, name, source, required)
	if err != nil || len(s) == 0 || !ok {
		return nil, false, err
	}

	return strings.Split(s, ","), true, nil
}

func GetInt32Array(r *http.Request, name string, required bool) ([]int32, bool, error) {
	pp, ok, err := GetStringArray(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	return convertInt32Array(pp)
}

func GetInt32ArrayFrom(r *http.Request, name string, source ParamSource, required bool) ([]int32, bool, error) {
	pp, ok, err := GetStringArrayFrom(r, name, source, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	return convertInt32Array(pp)
}

func convertInt32Array(pp []string) ([]int32, bool, error) {
	out := make([]int32, len(pp))

	for i, p := range pp {
		i32, err := strconv.ParseInt(p, 10, 32)
		if err != nil {
			return nil, false, fmt.Errorf("invalid int32:%q: %w", p, err)
		}

		out[i] = int32(i32)
	}

	return out, true, nil
}

func GetEnum[T comparable](r *http.Request, name string, required bool, parser func(string) T) (T, bool, error) {
	var out T

	s, ok, err := GetString(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return out, false, err
	}

	return parser(s), true, nil
}

func GetEnumFrom[T comparable](
	r *http.Request, name string, source ParamSource, required bool, parser func(string) T,
) (T, bool, error) {
	var out T

	s, ok, err := GetStringFrom(r, name, source, required)
	if err != nil || len(s) == 0 || !ok {
		return out, false, err
	}

	return parser(s), true, nil
}

func GetEnumArray[T comparable](
	r *http.Request, name string, required bool, parser func(string) T,
) ([]T, bool, error) {
	pp, ok, err := GetStringArray(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	out := make([]T, len(pp))

	for i, p := range pp {
		out[i] = parser(p)
	}

	return out, true, nil
}

func GetEnumArrayFrom[T comparable](
	r *http.Request, name string, source ParamSource, required bool, parser func(string) T,
) ([]T, bool, error) {
	pp, ok, err := GetStringArrayFrom(r, name, source, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	out := make([]T, len(pp))

	for i, p := range pp {
		out[i] = parser(p)
	}

	return out, true, nil
}

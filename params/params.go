package params

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

var ErrNotFound = errors.New("not found")

func GetString(r *http.Request, name string, required bool) (string, bool, error) {
	param := chi.URLParam(r, name)
	if len(param) == 0 {
		param = r.URL.Query().Get(name)
	}

	if required && len(param) == 0 {
		return "", false, fmt.Errorf("%s: %w", name, ErrNotFound)
	}

	return param, true, nil
}

func GetInt32(r *http.Request, name string, required bool) (int32, bool, error) {
	s, ok, err := GetString(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

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

	timestamp, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("invalid RFC3339 date: %w", err)
	}

	return timestamp, true, nil
}

func GetStringArray(r *http.Request, name string, required bool) ([]string, bool, error) {
	s, ok, err := GetString(r, name, required)
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

func GetEnumArray[T comparable](r *http.Request, name string, required bool, parser func(string) T) ([]T, bool, error) {
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

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

func GetString(r *http.Request, name string, required bool) (string, error) {
	param := chi.URLParam(r, name)
	if len(param) == 0 {
		param = r.URL.Query().Get(name)
	}

	if required && len(param) == 0 {
		return "", fmt.Errorf("%s: %w", name, ErrNotFound)
	}

	return param, nil
}

func GetInt32(r *http.Request, name string, required bool) (*int32, error) {
	s, err := GetString(r, name, required)
	if err != nil || len(s) == 0 {
		return nil, err
	}

	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid int32: %w", err)
	}

	i32 := int32(i)

	return &i32, nil
}

func GetInt(r *http.Request, name string, required bool) (*int, error) {
	s, err := GetString(r, name, required)
	if err != nil || len(s) == 0 {
		return nil, err
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("invalid int: %w", err)
	}

	return &i, nil
}

func GetTime(r *http.Request, name string, required bool) (time.Time, error) {
	s, err := GetString(r, name, required)
	if err != nil || len(s) == 0 {
		return time.Time{}, err
	}

	timestamp, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid RFC3339 date: %w", err)
	}

	return timestamp, nil
}

func GetStringArray(r *http.Request, name string, required bool) ([]string, error) {
	s, err := GetString(r, name, required)
	if err != nil || len(s) == 0 {
		return nil, err
	}

	return strings.Split(s, ","), nil
}

func GetInt32Array(r *http.Request, name string, required bool) ([]int32, error) {
	pp, err := GetStringArray(r, name, required)
	if err != nil || len(pp) == 0 {
		return nil, err
	}

	out := make([]int32, len(pp))

	for i, p := range pp {
		i32, err := strconv.ParseInt(p, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid int32:%q: %w", p, err)
		}

		out[i] = int32(i32)
	}

	return out, nil
}

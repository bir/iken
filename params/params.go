package params

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/bir/iken/validation"
)

const errFormat = "parameter invalid %s : %q : %w"

func GetString(r *http.Request, name string, required bool) (string, error) {
	param := chi.URLParam(r, name)
	if len(param) == 0 {
		param = r.URL.Query().Get(name)
	}

	if required && len(param) == 0 {
		return "", validation.New(name, name+" parameter not found") //nolint: wrapcheck
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
		return nil, validation.NewError(name, fmt.Errorf(errFormat, "int32", s, err)) //nolint: wrapcheck
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
		return nil, validation.NewError(name, fmt.Errorf(errFormat, "int", s, err)) //nolint: wrapcheck
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
		return time.Time{}, validation.NewError(name, fmt.Errorf(errFormat, "date", s, err)) //nolint: wrapcheck
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
			return nil, validation.NewError(name, fmt.Errorf(errFormat, "int32", p, err)) //nolint: wrapcheck
		}

		out[i] = int32(i32)
	}

	return out, nil
}

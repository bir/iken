package config

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/cast"
)

// ErrInvalidLocation is returned when a Location tag fails to load.
var ErrInvalidLocation = errors.New("failed parsing location")

// StringToLocationHookFunc converts strings to *time.Location.
func StringToLocationHookFunc(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	if t != reflect.TypeOf(time.Location{}) {
		return data, nil
	}

	s, ok := data.(string)
	if !ok {
		return data, nil
	}

	l, err := time.LoadLocation(s)
	if err != nil {
		return time.UTC, fmt.Errorf("%w: `%v`", ErrInvalidLocation, data)
	}

	return l, nil
}

// StringToMapStringStringHookFunc converts strings to map[string]string.
func StringToMapStringStringHookFunc(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	if t != reflect.TypeOf(map[string]string{}) {
		return data, nil
	}

	return cast.ToStringMapString(data), nil
}

// ErrInvalidURL is returned when a URL tag fails to parse.
var ErrInvalidURL = errors.New("failed parsing url")

// StringToURLHookFunc converts strings to *url.URL.
func StringToURLHookFunc(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	if t != reflect.TypeOf(&url.URL{}) { //nolint:exhaustruct
		return data, nil
	}

	s, ok := data.(string)
	if !ok {
		return data, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("%w: `%v`", ErrInvalidURL, data)
	}

	return u, nil
}

// ErrInvalidTime is returned when a time tag fails to parse.
var ErrInvalidTime = errors.New("failed parsing time")

// StringToTimeFunc converts strings to time.Time.
func StringToTimeFunc(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	if t != reflect.TypeOf(time.Time{}) {
		return data, nil
	}

	s, ok := data.(string)
	if !ok {
		return data, nil
	}

	out, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, fmt.Errorf("%w: `%v`", ErrInvalidTime, data)
	}

	return out, nil
}

// StringToRegexFunc converts strings to time.Time.
func StringToRegexFunc(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	if t != reflect.TypeOf(regexp.Regexp{}) {
		return data, nil
	}

	s, ok := data.(string)
	if !ok {
		return data, nil
	}

	out, err := regexp.Compile(s)
	if err != nil {
		return nil, fmt.Errorf("%w: `%v`", ErrInvalidTime, data)
	}

	return out, nil
}

// StringToSliceHookFunc returns a DecodeHookFunc that converts
// string to []string by splitting on the given sep.
func StringToSliceHookFunc(sep string) mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t.Kind() != reflect.Slice {
			return data, nil
		}

		raw, ok := data.(string)
		if !ok || raw == "" {
			return []string{}, nil
		}

		return strings.Split(raw, sep), nil
	}
}

func defaultDecoderConfig(c *mapstructure.DecoderConfig) {
	c.TagName = TagName
	c.DecodeHook = mapstructure.ComposeDecodeHookFunc(
		StringToLocationHookFunc,
		StringToMapStringStringHookFunc,
		StringToURLHookFunc,
		StringToTimeFunc,
		StringToRegexFunc,
		mapstructure.StringToTimeDurationHookFunc(),
		StringToSliceHookFunc(","),
	)
}

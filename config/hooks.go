package config

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"time"

	"github.com/mitchellh/mapstructure"
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

	l, err := time.LoadLocation(data.(string))
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

	u, err := url.Parse(data.(string))
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

	out, err := time.Parse(time.RFC3339, data.(string))
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

	out, err := regexp.Compile(data.(string))
	if err != nil {
		return nil, fmt.Errorf("%w: `%v`", ErrInvalidTime, data)
	}

	return out, nil
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
		mapstructure.StringToSliceHookFunc(","))
}

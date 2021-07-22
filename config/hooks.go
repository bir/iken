package config

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

// ErrInvalidLocation is returned when a Location tag fails to load.
var ErrInvalidLocation = errors.New("failed parsing location")

// StringToLocationHookFunc converts strings to *time.Location.
func StringToLocationHookFunc(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
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

func defaultDecoderConfig(c *mapstructure.DecoderConfig) {
	c.TagName = TagName
	c.DecodeHook = mapstructure.ComposeDecodeHookFunc(
		StringToLocationHookFunc,
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","))
}

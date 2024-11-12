package forms

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

type File struct {
	File     io.ReadCloser
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
}

func GetFile(r *http.Request, name string, required bool) (File, bool, error) {
	file, header, err := r.FormFile(name)
	if errors.Is(err, http.ErrMissingFile) {
		if required {
			return File{}, false, fmt.Errorf("%s: %w", name, ErrNotFound)
		}

		return File{}, false, nil
	}

	if err != nil {
		return File{}, false, fmt.Errorf("ToFormFile: %w", err)
	}

	return File{
		File:     file,
		Filename: header.Filename,
		Header:   header.Header,
		Size:     header.Size,
	}, true, nil
}

type LookupString func(key string) string

func GetString(lookup LookupString, name string, required bool) (string, error) {
	param := lookup(name)

	if required && len(param) == 0 {
		return "", fmt.Errorf("%s: %w", name, ErrNotFound)
	}

	return param, nil
}

func GetInt32(lookup LookupString, name string, required bool) (int32, error) {
	s, err := GetString(lookup, name, required)
	if err != nil || len(s) == 0 {
		return 0, err
	}

	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid int32: %w", err)
	}

	return int32(i), nil
}

func GetInt64(lookup LookupString, name string, required bool) (int64, error) {
	s, err := GetString(lookup, name, required)
	if err != nil || len(s) == 0 {
		return 0, err
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid int32: %w", err)
	}

	return i, nil
}

func GetBool(lookup LookupString, name string, required bool) (bool, error) {
	s, err := GetString(lookup, name, required)
	if err != nil || len(s) == 0 {
		return false, err
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, fmt.Errorf("invalid bool: %w", err)
	}

	return b, nil
}

func GetInt(lookup LookupString, name string, required bool) (int, error) {
	s, err := GetString(lookup, name, required)
	if err != nil || len(s) == 0 {
		return 0, err
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid int: %w", err)
	}

	return i, nil
}

func GetTime(lookup LookupString, name string, required bool) (time.Time, error) {
	s, err := GetString(lookup, name, required)
	if err != nil || len(s) == 0 {
		return time.Time{}, err
	}

	timestamp, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid RFC3339 date: %w", err)
	}

	return timestamp, nil
}

func GetUUID(lookup LookupString, name string, required bool) (uuid.UUID, error) {
	s, err := GetString(lookup, name, required)
	if err != nil || len(s) == 0 {
		return uuid.UUID{}, err
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid uuid: %w", err)
	}

	return id, nil
}

func GetStringArray(lookup LookupString, name string, required bool) ([]string, error) {
	s, err := GetString(lookup, name, required)
	if err != nil || len(s) == 0 {
		return nil, err
	}

	return strings.Split(s, ","), nil
}

func GetInt32Array(lookup LookupString, name string, required bool) ([]int32, error) {
	pp, err := GetStringArray(lookup, name, required)
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

func GetEnum[T comparable](lookup LookupString, name string, required bool, parser func(string) T) (T, error) {
	var out T

	s, err := GetString(lookup, name, required)
	if err != nil || len(s) == 0 {
		return out, err
	}

	return parser(s), nil
}

func GetEnumArray[T comparable](lookup LookupString, name string, required bool, parser func(string) T) ([]T, error) {
	pp, err := GetStringArray(lookup, name, required)
	if err != nil || len(pp) == 0 {
		return nil, err
	}

	out := make([]T, len(pp))

	for i, p := range pp {
		out[i] = parser(p)
	}

	return out, nil
}

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

func GetString(lookup LookupString, name string, required bool) (string, bool, error) {
	s := lookup(name)

	if s == "null" {
		s = ""
	}

	if required && len(s) == 0 {
		return "", false, fmt.Errorf("%s: %w", name, ErrNotFound)
	}

	return s, s != "", nil
}

func GetInt32(lookup LookupString, name string, required bool) (int32, bool, error) {
	s, ok, err := GetString(lookup, name, required)
	if err != nil || !ok {
		return 0, ok, err
	}

	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, false, fmt.Errorf("invalid int32: %w", err)
	}

	return int32(i), true, nil
}

func GetInt64(lookup LookupString, name string, required bool) (int64, bool, error) {
	s, ok, err := GetString(lookup, name, required)
	if err != nil || !ok {
		return 0, ok, err
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("invalid int32: %w", err)
	}

	return i, true, nil
}

func GetBool(lookup LookupString, name string, required bool) (bool, bool, error) {
	s, ok, err := GetString(lookup, name, required)
	if err != nil || !ok {
		return false, ok, err
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, false, fmt.Errorf("invalid bool: %w", err)
	}

	return b, true, nil
}

func GetInt(lookup LookupString, name string, required bool) (int, bool, error) {
	s, ok, err := GetString(lookup, name, required)
	if err != nil || !ok {
		return 0, ok, err
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false, fmt.Errorf("invalid int: %w", err)
	}

	return i, true, nil
}

func GetTime(lookup LookupString, name string, required bool) (time.Time, bool, error) {
	s, ok, err := GetString(lookup, name, required)
	if err != nil || !ok {
		return time.Time{}, ok, err
	}

	timestamp, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("invalid RFC3339 date: %w", err)
	}

	return timestamp, true, nil
}

func GetUUID(lookup LookupString, name string, required bool) (uuid.UUID, bool, error) {
	s, ok, err := GetString(lookup, name, required)
	if err != nil || !ok {
		return uuid.UUID{}, ok, err
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, false, fmt.Errorf("invalid uuid: %w", err)
	}

	return id, true, nil
}

func GetStringArray(lookup LookupString, name string, required bool) ([]string, bool, error) {
	s, ok, err := GetString(lookup, name, required)
	if err != nil || !ok {
		return nil, ok, err
	}

	return strings.Split(s, ","), true, nil
}

func GetInt32Array(lookup LookupString, name string, required bool) ([]int32, bool, error) {
	pp, ok, err := GetStringArray(lookup, name, required)
	if err != nil || len(pp) == 0 {
		return nil, ok, err
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

func GetUUIDArray(lookup LookupString, name string, required bool) ([]uuid.UUID, bool, error) {
	pp, ok, err := GetStringArray(lookup, name, required)
	if err != nil || len(pp) == 0 {
		return nil, ok, err
	}

	out := make([]uuid.UUID, len(pp))

	for i, p := range pp {
		id, err := uuid.Parse(p)
		if err != nil {
			return nil, false, fmt.Errorf("invalid uuid: %w", err)
		}

		out[i] = id
	}

	return out, true, nil
}

func GetEnum[T comparable](lookup LookupString, name string, required bool, parser func(string) T) (T, bool, error) {
	var out T

	s, ok, err := GetString(lookup, name, required)
	if err != nil || !ok {
		return out, ok, err
	}

	return parser(s), true, nil
}

func GetEnumArray[T comparable](fn LookupString, name string, required bool, parser func(string) T) ([]T, bool, error) {
	pp, ok, err := GetStringArray(fn, name, required)
	if err != nil || !ok {
		return nil, ok, err
	}

	out := make([]T, len(pp))

	for i, p := range pp {
		out[i] = parser(p)
	}

	return out, true, nil
}

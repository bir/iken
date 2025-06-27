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

var ErrNotFound = errors.New("not found")

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

func GetStringPath(r *http.Request, name string, required bool) (string, bool, error) {
	param := r.PathValue(name)
	if required && len(param) == 0 {
		return "", false, fmt.Errorf("%s: %w", name, ErrNotFound)
	}

	return param, param != "", nil
}

func GetStringQuery(r *http.Request, name string, required bool) (string, bool, error) {
	param := r.URL.Query().Get(name)
	if required && len(param) == 0 {
		return "", false, fmt.Errorf("%s: %w", name, ErrNotFound)
	}

	return param, param != "", nil
}

func GetStringHeader(r *http.Request, name string, required bool) (string, bool, error) {
	param := r.Header.Get(name)
	if required && len(param) == 0 {
		return "", false, fmt.Errorf("%s: %w", name, ErrNotFound)
	}

	return param, param != "", nil
}

func GetStringCookie(r *http.Request, name string, required bool) (string, bool, error) {
	cookie, err := r.Cookie(name)
	if err == nil && cookie != nil {
		// Return found even if cookie is empty, because it _is_ present!
		return cookie.Value, true, nil
	}

	if required {
		return "", false, fmt.Errorf("%s: %w", name, ErrNotFound)
	}

	return "", false, nil
}

func GetInt32(r *http.Request, name string, required bool) (int32, bool, error) {
	s, ok, err := GetString(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt32(s)
}

func GetInt32Path(r *http.Request, name string, required bool) (int32, bool, error) {
	s, ok, err := GetStringPath(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt32(s)
}

func GetInt32Query(r *http.Request, name string, required bool) (int32, bool, error) {
	s, ok, err := GetStringQuery(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt32(s)
}

func GetInt32Header(r *http.Request, name string, required bool) (int32, bool, error) {
	s, ok, err := GetStringHeader(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt32(s)
}

func GetInt32Cookie(r *http.Request, name string, required bool) (int32, bool, error) {
	s, ok, err := GetStringCookie(r, name, required)
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

func GetInt64Path(r *http.Request, name string, required bool) (int64, bool, error) {
	s, ok, err := GetStringPath(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt64(s)
}

func GetInt64Query(r *http.Request, name string, required bool) (int64, bool, error) {
	s, ok, err := GetStringQuery(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt64(s)
}

func GetInt64Header(r *http.Request, name string, required bool) (int64, bool, error) {
	s, ok, err := GetStringHeader(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt64(s)
}

func GetInt64Cookie(r *http.Request, name string, required bool) (int64, bool, error) {
	s, ok, err := GetStringCookie(r, name, required)
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

func GetBoolPath(r *http.Request, name string, required bool) (bool, bool, error) {
	s, ok, err := GetStringPath(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return false, false, err
	}

	return convertBool(s)
}

func GetBoolQuery(r *http.Request, name string, required bool) (bool, bool, error) {
	s, ok, err := GetStringQuery(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return false, false, err
	}

	return convertBool(s)
}

func GetBoolHeader(r *http.Request, name string, required bool) (bool, bool, error) {
	s, ok, err := GetStringHeader(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return false, false, err
	}

	return convertBool(s)
}

func GetBoolCookie(r *http.Request, name string, required bool) (bool, bool, error) {
	s, ok, err := GetStringCookie(r, name, required)
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

func GetIntPath(r *http.Request, name string, required bool) (int, bool, error) {
	s, ok, err := GetStringPath(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt(s)
}

func GetIntQuery(r *http.Request, name string, required bool) (int, bool, error) {
	s, ok, err := GetStringQuery(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt(s)
}

func GetIntHeader(r *http.Request, name string, required bool) (int, bool, error) {
	s, ok, err := GetStringHeader(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return 0, false, err
	}

	return convertInt(s)
}

func GetIntCookie(r *http.Request, name string, required bool) (int, bool, error) {
	s, ok, err := GetStringCookie(r, name, required)
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

func GetTimePath(r *http.Request, name string, required bool) (time.Time, bool, error) {
	s, ok, err := GetStringPath(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return time.Time{}, false, err
	}

	return convertTime(s)
}

func GetTimeQuery(r *http.Request, name string, required bool) (time.Time, bool, error) {
	s, ok, err := GetStringQuery(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return time.Time{}, false, err
	}

	return convertTime(s)
}

func GetTimeHeader(r *http.Request, name string, required bool) (time.Time, bool, error) {
	s, ok, err := GetStringHeader(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return time.Time{}, false, err
	}

	return convertTime(s)
}

func GetTimeCookie(r *http.Request, name string, required bool) (time.Time, bool, error) {
	s, ok, err := GetStringCookie(r, name, required)
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

func GetUUIDPath(r *http.Request, name string, required bool) (uuid.UUID, bool, error) {
	s, ok, err := GetStringPath(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return uuid.UUID{}, false, err
	}

	return convertUUID(s)
}

func GetUUIDQuery(r *http.Request, name string, required bool) (uuid.UUID, bool, error) {
	s, ok, err := GetStringQuery(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return uuid.UUID{}, false, err
	}

	return convertUUID(s)
}

func GetUUIDHeader(r *http.Request, name string, required bool) (uuid.UUID, bool, error) {
	s, ok, err := GetStringHeader(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return uuid.UUID{}, false, err
	}

	return convertUUID(s)
}

func GetUUIDCookie(r *http.Request, name string, required bool) (uuid.UUID, bool, error) {
	s, ok, err := GetStringCookie(r, name, required)
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

func GetStringArrayPath(r *http.Request, name string, required bool) ([]string, bool, error) {
	s, ok, err := GetStringPath(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return nil, false, err
	}

	return strings.Split(s, ","), true, nil
}

func GetStringArrayQuery(r *http.Request, name string, required bool) ([]string, bool, error) {
	s, ok, err := GetStringQuery(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return nil, false, err
	}

	return strings.Split(s, ","), true, nil
}

func GetStringArrayHeader(r *http.Request, name string, required bool) ([]string, bool, error) {
	s, ok, err := GetStringHeader(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return nil, false, err
	}

	return strings.Split(s, ","), true, nil
}

func GetStringArrayCookie(r *http.Request, name string, required bool) ([]string, bool, error) {
	s, ok, err := GetStringCookie(r, name, required)
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

func GetInt32ArrayPath(r *http.Request, name string, required bool) ([]int32, bool, error) {
	pp, ok, err := GetStringArrayPath(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	return convertInt32Array(pp)
}

func GetInt32ArrayQuery(r *http.Request, name string, required bool) ([]int32, bool, error) {
	pp, ok, err := GetStringArrayQuery(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	return convertInt32Array(pp)
}

func GetInt32ArrayHeader(r *http.Request, name string, required bool) ([]int32, bool, error) {
	pp, ok, err := GetStringArrayHeader(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	return convertInt32Array(pp)
}

func GetInt32ArrayCookie(r *http.Request, name string, required bool) ([]int32, bool, error) {
	pp, ok, err := GetStringArrayCookie(r, name, required)
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

func GetUUIDArray(r *http.Request, name string, required bool) ([]uuid.UUID, bool, error) {
	pp, ok, err := GetStringArray(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	return convertUUIDArray(pp)
}

func GetUUIDArrayPath(r *http.Request, name string, required bool) ([]uuid.UUID, bool, error) {
	pp, ok, err := GetStringArrayPath(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	return convertUUIDArray(pp)
}

func GetUUIDArrayQuery(r *http.Request, name string, required bool) ([]uuid.UUID, bool, error) {
	pp, ok, err := GetStringArrayQuery(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	return convertUUIDArray(pp)
}

func GetUUIDArrayHeader(r *http.Request, name string, required bool) ([]uuid.UUID, bool, error) {
	pp, ok, err := GetStringArrayHeader(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	return convertUUIDArray(pp)
}

func GetUUIDArrayCookie(r *http.Request, name string, required bool) ([]uuid.UUID, bool, error) {
	pp, ok, err := GetStringArrayCookie(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	return convertUUIDArray(pp)
}

func convertUUIDArray(pp []string) ([]uuid.UUID, bool, error) {
	out := make([]uuid.UUID, len(pp))

	for i, p := range pp {
		id, err := uuid.Parse(p)
		if err != nil {
			return nil, false, fmt.Errorf("invalid uuid:%q: %w", p, err)
		}

		out[i] = id
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

func GetEnumPath[T comparable](r *http.Request, name string, required bool, parser func(string) T) (T, bool, error) {
	var out T

	s, ok, err := GetStringPath(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return out, false, err
	}

	return parser(s), true, nil
}

func GetEnumQuery[T comparable](r *http.Request, name string, required bool, parser func(string) T) (T, bool, error) {
	var out T

	s, ok, err := GetStringQuery(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return out, false, err
	}

	return parser(s), true, nil
}

func GetEnumHeader[T comparable](r *http.Request, name string, required bool, parser func(string) T) (T, bool, error) {
	var out T

	s, ok, err := GetStringHeader(r, name, required)
	if err != nil || len(s) == 0 || !ok {
		return out, false, err
	}

	return parser(s), true, nil
}

func GetEnumCookie[T comparable](r *http.Request, name string, required bool, parser func(string) T) (T, bool, error) {
	var out T

	s, ok, err := GetStringCookie(r, name, required)
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

func GetEnumArrayPath[T comparable](
	r *http.Request, name string, required bool, parser func(string) T,
) ([]T, bool, error) {
	pp, ok, err := GetStringArrayPath(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	out := make([]T, len(pp))

	for i, p := range pp {
		out[i] = parser(p)
	}

	return out, true, nil
}

func GetEnumArrayQuery[T comparable](
	r *http.Request, name string, required bool, parser func(string) T,
) ([]T, bool, error) {
	pp, ok, err := GetStringArrayQuery(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	out := make([]T, len(pp))

	for i, p := range pp {
		out[i] = parser(p)
	}

	return out, true, nil
}

func GetEnumArrayHeader[T comparable](
	r *http.Request, name string, required bool, parser func(string) T,
) ([]T, bool, error) {
	pp, ok, err := GetStringArrayHeader(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	out := make([]T, len(pp))

	for i, p := range pp {
		out[i] = parser(p)
	}

	return out, true, nil
}

func GetEnumArrayCookie[T comparable](
	r *http.Request, name string, required bool, parser func(string) T,
) ([]T, bool, error) {
	pp, ok, err := GetStringArrayCookie(r, name, required)
	if err != nil || len(pp) == 0 || !ok {
		return nil, false, err
	}

	out := make([]T, len(pp))

	for i, p := range pp {
		out[i] = parser(p)
	}

	return out, true, nil
}

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

var ApplicationName string

// PgDBResolver is a wrapper to the GetPgDBString that adheres to the Resolver interface.
func PgDBResolver(_ reflect.StructField, key string) (any, bool, error) {
	value := GetPgDBString(key)

	return value, value != "", nil
}

func appendIf(aa []string, key, value string) []string {
	if value != "" {
		aa = append(aa, fmt.Sprintf("%s=%s", key, value))
	}

	return aa
}

// GetPgDBString reads env vars (via viper) based on the base name provide.
// for example GetPgDBString("DB") will read "DB_HOST", "DB_PORT", etc.  This is used
// in environments that manage the params separately.  Otherwise, just use a string type directly.
// application_name is attached to the connection string, it automatically uses the current running application.
// Set ApplicationName to override the value.
func GetPgDBString(base string) string {
	var pairs []string
	pairs = appendIf(pairs, "host", viper.GetString(base+"_HOST"))
	pairs = appendIf(pairs, "port", viper.GetString(base+"_PORT"))
	pairs = appendIf(pairs, "user", viper.GetString(base+"_USER"))
	pairs = appendIf(pairs, "password", viper.GetString(base+"_PASSWORD"))
	pairs = appendIf(pairs, "dbname", viper.GetString(base+"_NAME"))
	pairs = appendIf(pairs, "sslmode", viper.GetString(base+"_SSLMODE"))
	pairs = appendIf(pairs, "pool_max_conns", viper.GetString(base+"_MAX_CONN"))
	pairs = appendIf(pairs, "search_path", viper.GetString(base+"_SEARCH_PATH"))

	if len(pairs) == 0 {
		return ""
	}

	if ApplicationName == "" {
		ApplicationName = filepath.Base(os.Args[0])
	}

	pairs = append(pairs, "application_name="+ApplicationName)

	return strings.Join(pairs, " ")
}

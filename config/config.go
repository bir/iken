package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

var (
	// TagName defines the struct tag used to specify the env params.
	TagName = "env"
	// Type specifies the ConfigType used by viper.
	Type = "env"
	// File specifies the ConfigFile used by viper.  If the file does not exist then no error is raised.
	File = ".env"
	// Resolvers allow custom resolution for type mappings.  See Resolver for more info.
	Resolvers = ResolverMap{"pg": PgDBResolver}
	// ErrInvalidConfigObject is returned when a nil pointer, or non-pointer is provided to Load.
	ErrInvalidConfigObject = errors.New("config must be a pointer type")
	// ErrInvalidResolver is returned when a struct tag references a resolver that is not found.
	ErrInvalidResolver = errors.New("invalid resolver")
	// ErrInvalidTag is returned when a struct tag is improperly defined, e.g. `env:","`.
	ErrInvalidTag = errors.New("invalid tag")
)

// Resolver is used to map a key to a value.  Examples are custom serialization used for Postgres Connection URL
// composition, see PgDBResolver for example.
type Resolver func(key string) (interface{}, error)

// ResolverMap maintains the map of resolver names to Resolver funcs.
type ResolverMap = map[string]Resolver

const (
	keyPos      = 0
	defaultPos  = 1
	resolverPos = 2
)

// parseTag is responsible for parsing struct tags to env config.
// The struct tag format is:
//
//	`env:"ENV_VAR, DEFAULT(opt), RESOLVER(opt)"`
//
// Example Config struct:
//
//	   type ExampleConfig struct {
//		    DebugMode bool   `env:"DEBUG, false"`
//		    Port      int    `env:"PORT, 3000"`
//		    DB        string `env:"DB,localhost,pg"`
//	   }
func parseTag(tag string) error {
	args := strings.Split(tag, ",")

	key := strings.TrimSpace(args[keyPos])
	if key == "" {
		return fmt.Errorf("%w: `%s`", ErrInvalidTag, tag)
	}

	err := viper.BindEnv(key)
	if err != nil {
		return fmt.Errorf("binding tag: `%s`: %w", tag, err) // Ignore coverage - unlikely to error
	}

	if len(args) <= 1 {
		return nil
	}

	def := strings.TrimSpace(args[defaultPos])
	if def != "" {
		viper.SetDefault(key, def)
	}

	if len(args) == resolverPos {
		return nil
	}

	resolver := strings.TrimSpace(args[resolverPos])
	if resolver != "" {
		f, ok := Resolvers[resolver]
		if !ok {
			return fmt.Errorf("%w: `%v` for field `%v`", ErrInvalidResolver, resolver, key)
		}

		val, err := f(key)
		if err != nil {
			return err
		}

		viper.Set(key, val)
	}

	return nil
}

// Load uses struct tags (see parseTag) and viper to load the configuration into a config object.  The input
// object must be a pointer to a struct.  See ExampleLoad for simple example.
func Load(cfg interface{}) error {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.IsZero() {
		return ErrInvalidConfigObject
	}

	viper.SetConfigFile(File)
	viper.SetConfigType(Type)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	var pathError *os.PathError
	if err != nil && !errors.As(err, &pathError) {
		return fmt.Errorf("error loading config: %w", err)
	}

	v = reflect.Indirect(v)

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)

		tag := f.Tag.Get(TagName)
		if tag == "" || tag == "-" {
			continue
		}

		if err = parseTag(tag); err != nil {
			return fmt.Errorf("error parsing %s tag on field %s: %w", TagName, f.Name, err)
		}
	}

	err = viper.Unmarshal(cfg, defaultDecoderConfig)
	if err != nil {
		return fmt.Errorf("error Unmarshaling: %w", err)
	}

	return nil
}

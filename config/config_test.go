package config_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/bir/iken/config"
)

type Config struct {
	LocalDebug bool           `env:"DEBUG, false"`
	Port       int            `env:"PORT, 3000"`
	Interval   time.Duration  `env:"INTERVAL"`
	TimeZone   *time.Location `env:"TIMEZONE, America/Los_Angeles"`
	DB         string         `env:"DB,,pg"`
	MyUrl      *url.URL       `env:"MY_URL"`
	Ignore     string         `env:"-" json:"-"`
}

type InvalidConfig struct {
	BadTag string `env:","`
}

func ErrResolver(_ string) (any, error) {
	return nil, errors.New("BAD")
}

func TestSetup(t *testing.T) {
	// We use globals for some config settings, so these tests cannot run in parallel

	defaultResolvers := config.Resolvers
	defaultConfigFile := config.File
	b, _ := json.Marshal(Config{})
	defaultConfig := string(b)
	configWithService := fmt.Sprintf(`{"LocalDebug":true,"Port":1234,"Interval":0,"TimeZone":{},"DB":"host=1.2.3.4 user=user password=pass dbname=dbname application_name=%s","MyUrl":{"Scheme":"https","Opaque":"","User":null,"Host":"www.google.com","Path":"","RawPath":"","OmitHost":false,"ForceQuery":false,"RawQuery":"a=b","Fragment":"","RawFragment":""}}`, filepath.Base(os.Args[0]))

	tests := []struct {
		name    string
		pre     func()
		cfg     any
		env     map[string]string
		json    string
		wantErr bool
	}{
		{"defaults", nil, &Config{}, nil, `{"LocalDebug":true,"Port":1234,"Interval":0,"TimeZone":{},"DB":"host=1.2.3.4 user=user password=pass dbname=dbname application_name=test","MyUrl":{"Scheme":"https","Opaque":"","User":null,"Host":"www.google.com","Path":"","RawPath":"","OmitHost":false,"ForceQuery":false,"RawQuery":"a=b","Fragment":"","RawFragment":""}}`, false},
		{"envOverride", nil, &Config{}, map[string]string{"INTERVAL": "15s", "DB_PORT": "1", "DB_MAX_CONN": "99", "DB_SSLMODE": "funky"}, `{"LocalDebug":true,"Port":1234,"Interval":15000000000,"TimeZone":{},"DB":"host=1.2.3.4 port=1 user=user password=pass dbname=dbname sslmode=funky pool_max_conns=99 application_name=test","MyUrl":{"Scheme":"https","Opaque":"","User":null,"Host":"www.google.com","Path":"","RawPath":"","OmitHost":false,"ForceQuery":false,"RawQuery":"a=b","Fragment":"","RawFragment":""}}`, false},
		{"EmptyEnv", func() { config.File = ".envEMPTY" }, &Config{}, nil, `{"LocalDebug":false,"Port":3000,"Interval":0,"TimeZone":{},"DB":"","MyUrl":null}`, false},
		{"InvalidTZ", nil, &Config{}, map[string]string{"TIMEZONE": "FOO"}, `{"LocalDebug":true,"Port":1234,"Interval":0,"TimeZone":null,"DB":"host=1.2.3.4 user=user password=pass dbname=dbname application_name=test","MyUrl":{"Scheme":"https","Opaque":"","User":null,"Host":"www.google.com","Path":"","RawPath":"","OmitHost":false,"ForceQuery":false,"RawQuery":"a=b","Fragment":"","RawFragment":""}}`, true},
		{"InvalidURL", nil, &Config{}, map[string]string{"MY_URL": "%"}, `{"LocalDebug":true,"Port":1234,"Interval":0,"TimeZone":{},"DB":"host=1.2.3.4 user=user password=pass dbname=dbname application_name=test","MyUrl":null}`, true},
		{"BadConfig", nil, nil, nil, `null`, true},
		{"BadFile", func() { config.File = ".envBAD" }, &Config{}, nil, defaultConfig, true},
		{"BadResolver", func() { config.Resolvers = config.ResolverMap{} }, &Config{}, nil, defaultConfig, true},
		{"ErrResolver", func() { config.Resolvers = config.ResolverMap{"pg": ErrResolver} }, &Config{}, nil, defaultConfig, true},
		{"DefaultService", func() { config.ApplicationName = "" }, &Config{}, nil, configWithService, false},
		{"InvalidTag", nil, &InvalidConfig{}, nil, `{"BadTag":""}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.ApplicationName = "test"
			config.Resolvers = defaultResolvers
			config.File = defaultConfigFile
			if tt.pre != nil {
				tt.pre()
			}

			viper.Reset()
			os.Clearenv()
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			if err := config.Load(tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			b, err := json.Marshal(tt.cfg)
			if err != nil {
				t.Errorf("json.Marshall %v", err)
			}

			if string(b) != tt.json {
				t.Errorf("got \n%v, want \n%v", string(b), tt.json)
			}
		})
	}
}

type ComplexConfig struct {
	TestMap  map[string]string `env:"TEST_MAP"`
	Time     time.Time         `env:"TIME"`
	Patterns []regexp.Regexp   `json:",omitempty" env:"PATTERNS"`
}

func TestComplex(t *testing.T) {
	defaultResolvers := config.Resolvers
	defaultConfigFile := config.File

	tests := []struct {
		name    string
		pre     func()
		cfg     any
		env     map[string]string
		json    string
		wantErr bool
	}{
		{"defaults", nil, &ComplexConfig{}, nil, `{"TestMap":{"one":"1","two":"2"},"Time":"2021-01-01T00:00:00Z"}`, false},
		{"EmptyEnv", func() { config.File = ".envEMPTY" }, &ComplexConfig{}, nil, `{"TestMap":null,"Time":"0001-01-01T00:00:00Z"}`, false},
		{"Complex", func() { config.File = ".envCOMPLEX" }, &ComplexConfig{}, nil, `{"TestMap":{"one":"1","two":"2"},"Time":"2022-01-01T00:00:00Z","Patterns":["123","asdf","https://a.b/c","^http:*"]}`, false},
		{"BadTime", func() { config.File = ".envBADTIME" }, &ComplexConfig{}, nil, `{"TestMap":null,"Time":"0001-01-01T00:00:00Z"}`, true},
		{"BadMap", nil, &ComplexConfig{}, map[string]string{"TEST_MAP": "FOO"}, `{"TestMap":{},"Time":"2021-01-01T00:00:00Z"}`, false},
		{"BadRegex", func() { config.File = ".envEMPTY" }, &ComplexConfig{}, map[string]string{"PATTERNS": "a(b"}, `{"TestMap":null,"Time":"0001-01-01T00:00:00Z","Patterns":[""]}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.ApplicationName = "test"
			config.Resolvers = defaultResolvers
			config.File = defaultConfigFile
			if tt.pre != nil {
				tt.pre()
			}

			viper.Reset()
			os.Clearenv()
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			if err := config.Load(tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			b, err := json.Marshal(tt.cfg)
			if err != nil {
				t.Errorf("json.Marshall %v", err)
			}

			if string(b) != tt.json {
				t.Errorf("got \n%v, want \n%v", string(b), tt.json)
			}
		})
	}
}

type ExampleConfig struct {
	DebugMode bool     `env:"DEBUG, false"`
	Port      int      `env:"PORT, 3000"`
	DB        string   `env:"DB,localhost,pg"`
	Test      []string `env:"TEST_ARRAY"`
}

func ExampleLoad() {
	cfg := ExampleConfig{}
	_ = config.Load(&cfg)
	fmt.Printf("DebugMode=%v\n", cfg.DebugMode)
	fmt.Printf("Port=%v\n", cfg.Port)
	fmt.Printf("DB=%v\n", cfg.DB)
}

func TestFoo(t *testing.T) {
	t.Setenv("TEST_ARRAY", "1 2 3")
	t.Setenv("TEST_ARRAY2", "1 2 3")

	cfg := ExampleConfig{}

	_ = config.Load(&cfg)
	test := viper.GetStringSlice("TEST_ARRAY")
	fmt.Printf("TEST=%#v\n", cfg.Test)
	fmt.Printf("test=%#v\n", test)
}

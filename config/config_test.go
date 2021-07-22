package config_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bir/iken/config"
	"github.com/spf13/viper"
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

func ErrResolver(_ string) (interface{}, error) {
	return nil, errors.New("BAD")
}

func TestSetup(t *testing.T) {
	// We use globals for some config settings, so these tests cannot run in parallel
	defaultResolvers := config.Resolvers
	defaultConfigFile := config.File
	b, _ := json.Marshal(Config{})
	defaultConfig := string(b)
	configWithService := fmt.Sprintf(`{"LocalDebug":true,"Port":1234,"Interval":0,"TimeZone":{},"DB":"host=1.2.3.4 user=user password=pass dbname=dbname application_name=%s","MyUrl":{"Scheme":"https","Opaque":"","User":null,"Host":"www.google.com","Path":"","RawPath":"","ForceQuery":false,"RawQuery":"a=b","Fragment":""}}`, filepath.Base(os.Args[0]))

	tests := []struct {
		name    string
		pre     func()
		cfg     *Config
		env     map[string]string
		json    string
		wantErr bool
	}{
		{"defaults", nil, &Config{}, nil, `{"LocalDebug":true,"Port":1234,"Interval":0,"TimeZone":{},"DB":"host=1.2.3.4 user=user password=pass dbname=dbname application_name=test","MyUrl":{"Scheme":"https","Opaque":"","User":null,"Host":"www.google.com","Path":"","RawPath":"","ForceQuery":false,"RawQuery":"a=b","Fragment":""}}`, false},
		{"envOverride", nil, &Config{}, map[string]string{"INTERVAL": "15s", "DB_PORT": "1", "DB_MAX_CONN": "99", "DB_SSLMODE": "funky"}, `{"LocalDebug":true,"Port":1234,"Interval":15000000000,"TimeZone":{},"DB":"host=1.2.3.4 port=1 user=user password=pass dbname=dbname sslmode=funky pool_max_conns=99 application_name=test","MyUrl":{"Scheme":"https","Opaque":"","User":null,"Host":"www.google.com","Path":"","RawPath":"","ForceQuery":false,"RawQuery":"a=b","Fragment":""}}`, false},
		{"EmptyEnv", func() { config.File = ".envEMPTY" }, &Config{}, nil, `{"LocalDebug":false,"Port":3000,"Interval":0,"TimeZone":{},"DB":"","MyUrl":null}`, false},
		{"InvalidTZ", nil, &Config{}, map[string]string{"TIMEZONE": "FOO"}, `{"LocalDebug":true,"Port":1234,"Interval":0,"TimeZone":null,"DB":"host=1.2.3.4 user=user password=pass dbname=dbname application_name=test","MyUrl":{"Scheme":"https","Opaque":"","User":null,"Host":"www.google.com","Path":"","RawPath":"","ForceQuery":false,"RawQuery":"a=b","Fragment":""}}`, true},
		{"InvalidURL", nil, &Config{}, map[string]string{"MY_URL": "%"}, `{"LocalDebug":true,"Port":1234,"Interval":0,"TimeZone":{},"DB":"host=1.2.3.4 user=user password=pass dbname=dbname application_name=test","MyUrl":null}`, true},
		{"BadConfig", nil, nil, nil, `null`, true},
		{"BadFile", func() { config.File = ".envBAD" }, &Config{}, nil, defaultConfig, true},
		{"BadResolver", func() { config.Resolvers = config.ResolverMap{} }, &Config{}, nil, defaultConfig, true},
		{"ErrResolver", func() { config.Resolvers = config.ResolverMap{"pg": ErrResolver} }, &Config{}, nil, defaultConfig, true},
		{"DefaultService", func() { config.ApplicationName = "" }, &Config{}, nil, configWithService, false},
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
			for k, v := range tt.env {
				viper.Set(k, v)
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
	DebugMode bool   `env:"DEBUG, false"`
	Port      int    `env:"PORT, 3000"`
	DB        string `env:"DB,localhost,pg"`
}

func ExampleLoad() {
	cfg := ExampleConfig{}
	config.Load(&cfg)
	fmt.Printf("DebugMode=%v\n", cfg.DebugMode)
	fmt.Printf("Port=%v\n", cfg.Port)
	fmt.Printf("DB=%v\n", cfg.DB)
}

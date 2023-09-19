package config

import (
	"fmt"
	"github.com/spf13/viper"
	defstruct "marketplace-svc/pkg/default_struct"
	"os"
	"strings"
	"sync"
)

var m = &sync.Mutex{}
var c *Config

type KDOption struct {
	ConfigPath     string
	EnvPrefix      string
	ConfigFilename *string // filename without extension
}

type Option func(f *KDOption)

func WithConfigPath(configPath string) Option {
	return func(f *KDOption) {
		f.ConfigPath = configPath
	}
}

func WithConfigFilename(filename *string) Option {
	return func(f *KDOption) {
		f.ConfigFilename = filename
	}
}

func Init(opts ...Option) *Config {
	kdOption := &KDOption{ConfigPath: ".", EnvPrefix: "KD"}
	for _, applyOpt := range opts {
		applyOpt(kdOption)
	}

	k := NewViper(kdOption)
	err := k.Unmarshal(&c)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if e := defstruct.Set(c); e != nil {
		panic(fmt.Errorf("Fatal error when set the default config: %s \n", err))
	}

	return c
}

func Get(opts ...Option) *Config {
	if c == nil {
		m.Lock()
		defer m.Unlock()
		return Init(opts...)
	}
	return c
}

func NewViper(opt *KDOption) *viper.Viper {
	var configFileName string

	if opt.ConfigFilename != nil {
		// custom configFileName for testing
		configFileName = *opt.ConfigFilename
	} else {
		profile := "prd"
		if os.Getenv("KD_ENV") == "dev" || os.Getenv("KD_ENV") == "stg" {
			profile = "dev"
		}
		configFileName = configFileName + "config-" + profile
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(configFileName)
	v.AddConfigPath(opt.ConfigPath)
	v.SetEnvPrefix(opt.EnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return v
}

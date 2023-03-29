package viper

import (
	"bytes"
	"github.com/spf13/viper"
	"github.com/zedisdog/ty/config"
	"github.com/zedisdog/ty/errx"
	"strings"
)

func LoadEnv(enable bool) func(*Config) {
	return func(c *Config) {
		c.loadEnv = enable
	}
}

// NewConfig new config based viper.
//
//	 default options:
//			LoadEnv: true	//whether autoload environment variable
func NewConfig(opts ...func(*Config)) *Config {
	c := &Config{
		v:       viper.New(),
		loadEnv: true,
	}

	for _, set := range opts {
		set(c)
	}

	return c
}

type Config struct {
	v       *viper.Viper
	config  string
	loadEnv bool
}

func (c *Config) IsSet(key string) bool {
	return c.v.IsSet(key)
}

func (c *Config) Get(key string, def ...interface{}) interface{} {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.v.Get(key)
}

func (c *Config) GetString(key string, def ...string) string {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.v.GetString(key)
}

func (c *Config) GetInt(key string, def ...int) int {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.v.GetInt(key)
}

func (c *Config) GetBool(key string, def ...bool) bool {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.v.GetBool(key)
}

func (c *Config) GetIntSlice(key string, def ...[]int) []int {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.v.GetIntSlice(key)
}

func (c *Config) GetStringSlice(key string, def ...[]string) []string {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.v.GetStringSlice(key)
}

func (c *Config) GetStringMap(key string, def ...map[string]interface{}) map[string]interface{} {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.v.GetStringMap(key)
}

func (c *Config) Sub(key string) config.IConfig {
	return &Config{
		v: c.v.Sub(key),
	}
}

func (c *Config) Set(key string, value any) {
	c.v.Set(key, value)
}

func (c *Config) SetConfigType(typeStr string) config.IConfig {
	c.v.SetConfigType(typeStr)
	return c
}

func (c *Config) SetConfig(config string) config.IConfig {
	c.config = config
	return c
}

func (c *Config) SetEnvKeyReplacer(replacer *strings.Replacer) config.IConfig {
	c.v.SetEnvKeyReplacer(replacer)
	return c
}

func (c *Config) Load() {
	err := c.v.ReadConfig(bytes.NewBuffer([]byte(c.config)))
	if err != nil {
		panic(errx.Wrap(err, "[config]read config failed"))
	}
	c.v.AutomaticEnv()
}

func (c *Config) AllSettings() interface{} {
	return c.v.AllSettings()
}

package viper

import (
	"bytes"
	"github.com/spf13/viper"
	"github.com/zedisdog/ty/config"
	"github.com/zedisdog/ty/errx"
	"io"
	"strings"
)

func NewConfig() config.IConfig {
	return &Config{
		v: viper.New(),
	}
}

type Config struct {
	v      *viper.Viper
	config []byte
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

func (c *Config) SetYml(config []byte) {
	c.v.SetConfigType("yml")
	c.config = config
}

func (c *Config) SetEnvKeyReplacer(replacer *strings.Replacer) {
	c.v.SetEnvKeyReplacer(replacer)
}

func (c *Config) Load() {
	err := c.v.ReadConfig(bytes.NewBuffer(c.config))
	if err != nil {
		panic(errx.Wrap(err, "[config]read config failed"))
	}
	c.v.AutomaticEnv()
}

func (c *Config) New(cfg interface{}) (conf config.IConfig, err error) {
	v := viper.New()
	switch c := cfg.(type) {
	case io.Reader:
		err = v.MergeConfig(c)
	case map[string]interface{}:
		err = v.MergeConfigMap(c)
	default:
		err = errx.New("config is invalid")
	}

	if err != nil {
		return
	}

	conf = &Config{
		v: v,
	}

	return
}

func (c *Config) AllSettings() interface{} {
	return c.v.AllSettings()
}

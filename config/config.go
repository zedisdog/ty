package config

import (
	"bytes"
	"github.com/spf13/viper"
	"strings"
)

func NewConfigWithBytesContent(configType string, content []byte, opts ...func(*Config)) *Config {
	c := NewConfig(opts...)
	c.SetConfigType(configType)
	err := c.LoadBytes(content)
	if err != nil {
		panic(err)
	}
	c.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	c.AutomaticEnv()

	return c
}

// NewConfig new config based viper.
func NewConfig(opts ...func(*Config)) *Config {
	c := &Config{
		Viper: viper.New(),
	}

	for _, set := range opts {
		set(c)
	}

	return c
}

type Config struct {
	*viper.Viper
}

func (c *Config) Get(key string, def ...interface{}) interface{} {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.Viper.Get(key)
}

func (c *Config) GetString(key string, def ...string) string {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.Viper.GetString(key)
}

func (c *Config) GetInt(key string, def ...int) int {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.Viper.GetInt(key)
}

func (c *Config) GetBool(key string, def ...bool) bool {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.Viper.GetBool(key)
}

func (c *Config) GetIntSlice(key string, def ...[]int) []int {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.Viper.GetIntSlice(key)
}

func (c *Config) GetStringSlice(key string, def ...[]string) []string {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.Viper.GetStringSlice(key)
}

func (c *Config) GetStringMap(key string, def ...map[string]interface{}) map[string]interface{} {
	if !c.IsSet(key) && len(def) > 0 {
		return def[0]
	}

	return c.Viper.GetStringMap(key)
}

func (c *Config) Sub(key string) *Config {
	v := c.Viper.Sub(key)
	if v == nil {
		return nil
	}

	return &Config{
		Viper: v,
	}
}

func (c *Config) SetConfigType(typeStr string) *Config {
	c.Viper.SetConfigType(typeStr)
	return c
}

func (c *Config) SetEnvKeyReplacer(replacer *strings.Replacer) *Config {
	c.Viper.SetEnvKeyReplacer(replacer)
	return c
}

func (c *Config) LoadBytes(content []byte) error {
	return c.Viper.ReadConfig(bytes.NewBuffer(content))
}

func (c *Config) LoadString(content string) error {
	return c.LoadBytes([]byte(content))
}

func (c *Config) LoadEnv() {
	c.Viper.AutomaticEnv()
}

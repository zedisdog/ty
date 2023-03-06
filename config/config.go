package config

import (
	"strings"
)

type IConfig interface {
	SetYml(config []byte)
	SetEnvKeyReplacer(replacer *strings.Replacer)
	// Load loads all config. panic if there is error.
	Load()
	Get(key string, def ...interface{}) interface{}
	GetString(key string, def ...string) string
	GetInt(key string, def ...int) int
	GetBool(key string, def ...bool) bool
	IsSet(key string) bool
	Sub(key string) IConfig
	GetIntSlice(key string, def ...[]int) []int
	GetStringSlice(key string, def ...[]string) []string
	GetStringMap(key string, def ...map[string]interface{}) map[string]interface{}
}

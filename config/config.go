package config

import (
	"strings"
)

type CanSetEnvKeyReplacer interface {
	SetEnvKeyReplacer(replacer *strings.Replacer) IConfig
}

type CanSetConfigType interface {
	SetConfigType(typeStr string) IConfig
}

type TypeGets interface {
	GetString(key string, def ...string) string
	GetInt(key string, def ...int) int
	GetBool(key string, def ...bool) bool
	GetIntSlice(key string, def ...[]int) []int
	GetStringSlice(key string, def ...[]string) []string
	GetStringMap(key string, def ...map[string]interface{}) map[string]interface{}
}

type IConfig interface {
	SetConfig(pathOrContent string) IConfig

	Set(key string, value interface{})

	// Load loads all config. panic if there is error.
	Load()
	Sub(key string) IConfig
	IsSet(key string) bool

	AllSettings() interface{}
	Get(key string, def ...interface{}) interface{}

	TypeGets
	CanSetConfigType
	CanSetEnvKeyReplacer
}

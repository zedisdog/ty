package application

import "github.com/zedisdog/ty/config"

type IModule interface {
	Name() string
	Register() error
	Bootstrap(config config.IConfig) error
}

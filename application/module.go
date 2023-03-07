package application

import "github.com/zedisdog/ty/config"

type IModule interface {
	Name() string
	Register(application IApplication) error
	Bootstrap(config config.IConfig) error
}

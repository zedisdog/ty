package application

import "github.com/zedisdog/ty/config"

type IModule interface {
	Name() string
	// Register registers resource to application. e.g: route used by default http server
	Register(application IApplication) error
	// Bootstrap start module's own sub process.
	Bootstrap(config config.IConfig) error
}

package application

import (
	"github.com/zedisdog/ty/config"
)

type IHasConfig interface {
	Config() config.IConfig
}

func SetConfig(config *config.Config) {
	GetInstance().SetConfig(config)
}
func (app *App) SetConfig(config config.IConfig) {
	app.config = config
}

func Config() config.IConfig {
	return GetInstance().Config()
}
func (app *App) Config() config.IConfig {
	return app.config
}

package application

import (
	"github.com/zedisdog/ty/config"
)

type IHasConfig interface {
	Config() *config.Config
}

func SetConfig(config *config.Config) {
	GetInstance().SetConfig(config)
}
func (app *App) SetConfig(config *config.Config) {
	app.config = config
}

func Config() *config.Config {
	return GetInstance().Config()
}
func (app *App) Config() *config.Config {
	return app.config
}

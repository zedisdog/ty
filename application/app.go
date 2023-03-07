package application

import (
	"fmt"
	"github.com/zedisdog/ty/config"
	"github.com/zedisdog/ty/errx"
	"github.com/zedisdog/ty/log"
	"github.com/zedisdog/ty/log/zap"
	"github.com/zedisdog/ty/server"
	"sync"
)

var Instance *App

func InitApp(config config.IConfig) {
	Instance = &App{
		config:      config,
		httpServers: new(sync.Map),
	}
}

type IApplication interface {
}

type App struct {
	config      config.IConfig
	httpServers *sync.Map
	logger      log.ILog
	modules     []IModule
}

func (app *App) Bootstrap() {
	app.initLog(app.config.Sub("log"))
	app.bootModules(app.config.Sub("modules"))
}

func (app *App) initLog(config config.IConfig) {
	driver := config.GetString("driver", "zap")
	fmt.Printf("[bootstrap]init log using %s...", driver)
	switch driver {
	case "zap":
		app.logger = zap.NewZapLog()
	}
	return
}

func (app *App) bootModules(config config.IConfig) {
	for _, module := range app.modules {
		if config.IsSet(module.Name()) {
			err := module.Bootstrap(config.Sub(module.Name()))
			if err != nil {
				panic(errx.Wrap(err, "[app]boot module failed"))
			}
		}
	}
}

func RegisterHttpServer(name string, server server.IHTTPServer) {
	Instance.RegisterHttpServer(name, server)
}
func (app *App) RegisterHttpServer(name string, server server.IHTTPServer) {
	app.httpServers.Store(name, server)
}

func GetHttpServer(name string) server.IHTTPServer {
	return Instance.GetHttpServer(name)
}
func (app *App) GetHttpServer(name string) server.IHTTPServer {
	v, exists := app.httpServers.Load(name)
	if !exists {
		return nil
	}

	return v.(server.IHTTPServer)
}

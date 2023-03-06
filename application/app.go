package application

import (
	"fmt"
	"github.com/zedisdog/ty/config"
	"github.com/zedisdog/ty/log"
	"github.com/zedisdog/ty/log/zap"
	"github.com/zedisdog/ty/server"
	"sync"
)

var Instance *App

func InitApp(config config.IConfig) {
	Instance = &App{
		config:     config,
		httpServer: new(sync.Map),
	}
}

type App struct {
	config     config.IConfig
	httpServer *sync.Map
	logger     log.ILog
}

func (app *App) Bootstrap() {
	app.initLog(app.config.Sub("log"))
	app.initModules(app.config.GetStringMap("modules"))
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

func (app *App) initModules(config map[string]interface{}) {

}

func RegisterHttpServer(name string, server server.IHTTPServer) {
	Instance.RegisterHttpServer(name, server)
}
func (app *App) RegisterHttpServer(name string, server server.IHTTPServer) {
	app.httpServer.Store(name, server)
}

func GetHttpServer(name string) server.IHTTPServer {
	return Instance.GetHttpServer(name)
}
func (app *App) GetHttpServer(name string) server.IHTTPServer {
	v, exists := app.httpServer.Load(name)
	if !exists {
		return nil
	}

	return v.(server.IHTTPServer)
}

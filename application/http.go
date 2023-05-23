package application

import (
	"fmt"
	"github.com/zedisdog/ty/sdk/net/http/server"
	"github.com/zedisdog/ty/sdk/net/http/server/gin"
)

type IHasHttpServer interface {
	RegisterHttpServer(name string, svr any)
	HttpServer(name ...string) any
}

func RegisterHttpServer[T any](name string, svr server.IHTTPServer[T]) {
	GetInstance().RegisterHttpServer(name, svr)
}
func (app *App) RegisterHttpServer(name string, svr any) {
	app.httpServers.Store(name, svr)
}

func HttpServer[T any](name ...string) server.IHTTPServer[T] {
	return GetInstance().HttpServer(name...).(server.IHTTPServer[T])
}
func (app *App) HttpServer(name ...string) any {
	return app.getValueOrDefault(app.httpServers, name...)
}

func (app *App) initDefaultHttpServer() {
	config := app.config.Sub("default.httpServer")
	if config != nil {
		app.logger.Info("[application] create default http server...")
		svr := gin.NewGinServer(
			fmt.Sprintf(
				"%s:%d",
				config.GetString("host"),
				config.GetInt("port"),
			),
			gin.EnablePprof(config.GetBool("enablePprof")),
			gin.WithLogger(app.logger),
		)
		app.httpServers.Store("default", svr)
	}
}

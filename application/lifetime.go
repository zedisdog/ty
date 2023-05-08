package application

import (
	"fmt"
	"github.com/zedisdog/ty/config"
	"github.com/zedisdog/ty/log"
	"os"
	"os/signal"
	"syscall"
)

type ILifetime interface {
	Init(config *config.Config)
	RegisterSeeder(seeders ...func() error)
	RegisterStopFunc(f func())
	Boot()
	Run()
	Wait(closeFunc ...func())
}

/*********************************init*********************************************/

// Init set config to application.
func Init(config *config.Config) {
	GetInstance().Init(config)
}
func (app *App) Init(config *config.Config) {
	app.SetConfig(config)

	app.initLog()
	app.initDefaultDatabase()
	app.initDefaultMigrator()
	app.initDefaultStorage()
	app.initDefaultHttpServer()
	app.RegisterStopFunc(app.stop)
}

/****************************************register***********************************/

func RegisterStopFunc(f func()) {
	GetInstance().RegisterStopFunc(f)
}
func (app *App) RegisterStopFunc(f func()) {
	app.onStop = append([]func(){f}, app.onStop...)
}

/******************************************run**********************************************/

// Boot boots the application.
func Boot() {
	GetInstance().Boot()
}
func (app *App) Boot() {
	app.migrate()
	app.seed()
	app.modules.Range(func(key, module any) bool {
		app.logger.Info(fmt.Sprintf("boot module <%s>...", module.(IModule).Name()))
		err := module.(IModule).Boot()
		if err != nil {
			panic(err)
		}
		return true
	})
}

func Run() {
	GetInstance().Run()
}
func (app *App) Run() {
	app.httpServers.Range(func(key, value any) bool {
		app.logger.Info(
			"[application] run http server...",
			&log.Field{Name: "name", Value: key},
		)
		value.(interface{ Run() }).Run()
		return true
	})
}

func (app *App) stop() {
	app.httpServers.Range(func(key, value any) bool {
		app.logger.Info(
			"[application] shutdown http server...",
			&log.Field{Name: "name", Value: key},
		)
		err := value.(interface{ Shutdown() error }).Shutdown()
		if err != nil {
			app.logger.Error(
				"[application]shutdown http server error",
				&log.Field{Name: "name", Value: key},
				&log.Field{Name: "error", Value: err},
			)
		}
		return true
	})
	app.CloseScheduler()
}

func Wait(closeFunc ...func()) {
	GetInstance().Wait(closeFunc...)
}
func (app *App) Wait(closeFunc ...func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-c
	for _, cls := range append(app.onStop, closeFunc...) {
		cls()
	}
}
